package repository

import (
	"context"
	"errors"

	"github.com/krisnadwipayana07/restful-fintech/internal/domain/model"
	"gorm.io/gorm"
)

type WalletRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Wallet, error)
	CreateWallet(ctx context.Context, account *model.Wallet) error
	UpdateBalance(ctx context.Context, tx *gorm.DB, walletID int64, newBalance string) error
}

type WalletRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) WalletRepository {
	return &WalletRepositoryImpl{db: db}
}

func (r *WalletRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Wallet, error) {
	var account model.Wallet
	err := r.db.WithContext(ctx).Take(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Explicitly return nil to indicate no data found
		}
		return nil, err
	}
	return &account, nil
}

func (r *WalletRepositoryImpl) CreateWallet(ctx context.Context, account *model.Wallet) error {
	return r.db.WithContext(ctx).
		Create(account).
		Error
}

func (r *WalletRepositoryImpl) UpdateBalance(ctx context.Context, tx *gorm.DB, walletID int64, newBalance string) error {
	return tx.WithContext(ctx).
		Model(&model.Wallet{}).
		Where("id = ?", walletID).
		Update("wallet_curr_balance", newBalance).
		Error
}
