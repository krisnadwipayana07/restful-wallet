package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/krisnadwipayana07/restful-fintech/internal/domain/constant"
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/model"
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/repository"
	"github.com/krisnadwipayana07/restful-fintech/pkg/dto"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionService interface {
	Withdraw(ctx context.Context, idempotencyKey string, walletID int64, req dto.AmountRequest) (resp dto.TransactionResponse, err error)
	Deposit(ctx context.Context, idempotencyKey string, walletID int64, req dto.AmountRequest) (resp dto.TransactionResponse, err error)
	Transfer(ctx context.Context, idempotencyKey string, walletID int64, req dto.TransferRequest) (resp dto.TransactionResponse, err error)
}

type TransactionServiceImpl struct {
	db              *gorm.DB
	redis           *redis.Client
	transactionRepo repository.TransactionRepository
	walletRepo      repository.WalletRepository
}

func NewTransactionService(db *gorm.DB, redis *redis.Client, repo repository.TransactionRepository, walletRepo repository.WalletRepository) TransactionService {
	return &TransactionServiceImpl{db: db, redis: redis, transactionRepo: repo, walletRepo: walletRepo}
}

func (s *TransactionServiceImpl) Withdraw(ctx context.Context, idempotencyKey string, walletID int64, req dto.AmountRequest) (resp dto.TransactionResponse, err error) {
	// Check double request
	val, err := s.redis.Exists(ctx, idempotencyKey).Result()
	if val == 1 {
		return dto.TransactionResponse{}, errors.New("Double Request")
	}

	// Set Idempotency Key
	s.redis.Set(ctx, idempotencyKey, true, 24*time.Hour).Err()

	curretWallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	currBalance, err := decimal.NewFromString(curretWallet.CurrentBalance)
	if err != nil {
		log.Printf("error converting existing wallet to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}
	wdAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error wd amount to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	if wdAmount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to withdraw 0 amount, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to 0 amount")
	}

	if wdAmount.GreaterThan(currBalance) {
		log.Printf("attempting to withdraw more than available balance, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to withdraw more than available balance")
	}

	// begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return dto.TransactionResponse{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// create transaction
	transactionID, err := s.transactionRepo.CreateTransaction(ctx, tx, model.Transaction{
		ID:        0,
		WalletID:  curretWallet.ID,
		Type:      constant.TransactionTypeWithdraw,
		IsDebit:   false,
		Value:     wdAmount.String(),
		Remarks:   "Withdraw",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("creating transaction, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	// update balance
	newBalance := currBalance.Sub(wdAmount)
	err = s.walletRepo.UpdateBalance(ctx, tx, curretWallet.ID, newBalance.String())
	if err != nil {
		log.Printf("updating new balance, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
		TransactionID: transactionID,
	}, nil
}

func (s *TransactionServiceImpl) Deposit(ctx context.Context, idempotencyKey string, walletID int64, req dto.AmountRequest) (resp dto.TransactionResponse, err error) {
	// Check double request
	val, err := s.redis.Exists(ctx, idempotencyKey).Result()
	if val == 1 {
		return dto.TransactionResponse{}, errors.New("Double Request")
	}
	// Set Idempotency Key
	s.redis.Set(ctx, idempotencyKey, true, 24*time.Hour).Err()

	depositAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error wd amount to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	if depositAmount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to deposit 0 or less, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to deposit 0 or less")
	}

	// begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return dto.TransactionResponse{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// create transaction
	transactionID, err := s.transactionRepo.CreateTransaction(ctx, tx, model.Transaction{
		ID:        0,
		WalletID:  walletID,
		Type:      constant.TransactionTypeDeposit,
		IsDebit:   true,
		Value:     depositAmount.String(),
		Remarks:   "Deposit",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("creating transaction, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	curretWallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	currBalance, err := decimal.NewFromString(curretWallet.CurrentBalance)
	if err != nil {
		log.Printf("error converting existing wallet to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}
	// update balance
	newBalance := currBalance.Add(depositAmount)
	err = s.walletRepo.UpdateBalance(ctx, tx, curretWallet.ID, newBalance.String())
	if err != nil {
		log.Printf("updating new balance, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
		TransactionID: transactionID,
	}, nil
}

func (s *TransactionServiceImpl) Transfer(ctx context.Context, idempotencyKey string, walletID int64, req dto.TransferRequest) (resp dto.TransactionResponse, err error) {
	// Check double request
	val, err := s.redis.Exists(ctx, idempotencyKey).Result()
	if val == 1 {
		return dto.TransactionResponse{}, errors.New("Double Request")
	}
	// Set Idempotency Key
	s.redis.Set(ctx, idempotencyKey, true, 24*time.Hour).Err()

	curretWallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	senderBalance, err := decimal.NewFromString(curretWallet.CurrentBalance)
	if err != nil {
		log.Printf("error converting existing wallet to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}
	transferAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		log.Printf("error wd amount to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	if transferAmount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to transfer 0 amount, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to 0 amount")
	}

	if transferAmount.GreaterThan(senderBalance) {
		log.Printf("attempting to transfer more than available balance, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to transfer more than available balance")
	}

	// begin transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return dto.TransactionResponse{}, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// create transaction
	senderTransactionID, err := s.transactionRepo.CreateTransaction(ctx, tx, model.Transaction{
		ID:        0,
		WalletID:  curretWallet.ID,
		Type:      constant.TransactionTypeTransfer,
		IsDebit:   false,
		Value:     transferAmount.String(),
		Remarks:   "Transfer - Send",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("creating transaction, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	// update balance
	newSenderBalance := senderBalance.Sub(transferAmount)
	err = s.walletRepo.UpdateBalance(ctx, tx, curretWallet.ID, newSenderBalance.String())
	if err != nil {
		log.Printf("updating new balance, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	// create transaction
	_, err = s.transactionRepo.CreateTransaction(ctx, tx, model.Transaction{
		ID:        0,
		WalletID:  req.ReceiverWalletID,
		Type:      constant.TransactionTypeTransfer,
		IsDebit:   true,
		Value:     transferAmount.String(),
		Remarks:   "Transfer - Receive",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("creating transaction, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	receiverWallet, err := s.walletRepo.FindByID(ctx, req.ReceiverWalletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	receiverBalance, err := decimal.NewFromString(receiverWallet.CurrentBalance)
	if err != nil {
		log.Printf("error converting existing wallet to decimal, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	// update balance
	newReceiverBalance := receiverBalance.Add(transferAmount)
	err = s.walletRepo.UpdateBalance(ctx, tx, req.ReceiverWalletID, newReceiverBalance.String())
	if err != nil {
		log.Printf("updating new balance, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	return dto.TransactionResponse{
			TransactionID: senderTransactionID,
		},
		nil
}
