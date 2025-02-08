package repository

import (
	"context"

	"github.com/krisnadwipayana07/restful-fintech/internal/domain/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Transaction, error)
	CreateTransaction(ctx context.Context, tx *gorm.DB, transaction model.Transaction) (int64, error)
	GetListTransactionByWalletID(ctx context.Context, walletID int64) ([]model.Transaction, error)
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepositoryImpl) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction model.Transaction) (int64, error) {
	err := tx.WithContext(ctx).Create(&transaction).Error
	return transaction.ID, err
}

func (r *TransactionRepositoryImpl) GetListTransactionByWalletID(ctx context.Context, walletID int64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Find(&transactions).Error
	return transactions, err
}
