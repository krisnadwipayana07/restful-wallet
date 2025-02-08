package service

import (
	"context"
	"errors"

	"github.com/krisnadwipayana07/restful-fintech/internal/domain/repository"
	"github.com/krisnadwipayana07/restful-fintech/pkg/dto"
)

type WalletService interface {
	WalletHistory(ctx context.Context, id int64) ([]dto.TransactionDetailResponse, error)
	WalletBalance(ctx context.Context, id int64) (string, error)
}

type WalletServiceImpl struct {
	transactionRepo repository.TransactionRepository
	walletRepo      repository.WalletRepository
}

func NewWalletService(repo repository.TransactionRepository, walletRepo repository.WalletRepository) WalletService {
	return &WalletServiceImpl{transactionRepo: repo, walletRepo: walletRepo}
}

func (s *WalletServiceImpl) WalletHistory(ctx context.Context, id int64) ([]dto.TransactionDetailResponse, error) {
	data, err := s.transactionRepo.GetListTransactionByWalletID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.NewTransactionListResponse(data), nil
}

func (s *WalletServiceImpl) WalletBalance(ctx context.Context, id int64) (string, error) {
	data, err := s.walletRepo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", errors.New("wallet not found")
	}
	return data.CurrentBalance, nil
}
