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

func (s *TransactionServiceImpl) createTransactionWithUpdateBalance(ctx context.Context, tx *gorm.DB, transaction model.Transaction, newBalance decimal.Decimal) (int64, error) {
	transactionID, err := s.transactionRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		log.Printf("creating transaction, err: %+v", err)
		return 0, err
	}

	err = s.walletRepo.UpdateBalance(ctx, tx, transaction.WalletID, newBalance)
	if err != nil {
		log.Printf("updating new balance, err: %+v", err)
		return 0, err
	}

	return transactionID, nil
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

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to withdraw 0 amount, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to 0 amount")
	}

	if req.Amount.GreaterThan(curretWallet.CurrentBalance) {
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
	transaction := model.Transaction{
		ID:        0,
		WalletID:  curretWallet.ID,
		Type:      constant.TransactionTypeWithdraw,
		IsDebit:   false,
		Value:     req.Amount,
		Remarks:   "Withdraw",
		CreatedAt: time.Now(),
	}
	newBalance := curretWallet.CurrentBalance.Sub(req.Amount)
	transactionID, err := s.createTransactionWithUpdateBalance(ctx, tx, transaction, newBalance)
	if err != nil {
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

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to deposit 0 or less, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to deposit 0 or less")
	}

	curretWallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
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

	transaction := model.Transaction{
		ID:        0,
		WalletID:  walletID,
		Type:      constant.TransactionTypeDeposit,
		IsDebit:   true,
		Value:     req.Amount,
		Remarks:   "Deposit",
		CreatedAt: time.Now(),
	}
	newBalance := curretWallet.CurrentBalance.Add(req.Amount)

	transactionID, err := s.createTransactionWithUpdateBalance(ctx, tx, transaction, newBalance)
	if err != nil {
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

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		log.Printf("attempting to transfer 0 amount, err: %+v", err)
		return dto.TransactionResponse{}, errors.New("attempting to 0 amount")
	}

	if req.Amount.GreaterThan(curretWallet.CurrentBalance) {
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
	senderTransaction := model.Transaction{
		ID:        0,
		WalletID:  curretWallet.ID,
		Type:      constant.TransactionTypeTransfer,
		IsDebit:   false,
		Value:     req.Amount,
		Remarks:   "Transfer - Send",
		CreatedAt: time.Now(),
	}
	newSenderBalance := curretWallet.CurrentBalance.Sub(req.Amount)
	senderTransactionID, err := s.createTransactionWithUpdateBalance(ctx, tx, senderTransaction, newSenderBalance)
	if err != nil {
		return dto.TransactionResponse{}, err
	}

	// create transaction
	receiverTransaction := model.Transaction{
		ID:        0,
		WalletID:  req.ReceiverWalletID,
		Type:      constant.TransactionTypeTransfer,
		IsDebit:   true,
		Value:     req.Amount,
		Remarks:   "Transfer - Receive",
		CreatedAt: time.Now(),
	}

	receiverWallet, err := s.walletRepo.FindByID(ctx, req.ReceiverWalletID)
	if err != nil {
		log.Printf("error wallet find by id, err: %+v", err)
		return dto.TransactionResponse{}, err
	}

	newReceiverBalance := receiverWallet.CurrentBalance.Add(req.Amount)
	_, err = s.createTransactionWithUpdateBalance(ctx, tx, receiverTransaction, newReceiverBalance)
	if err != nil {
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
