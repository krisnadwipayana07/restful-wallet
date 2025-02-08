package repository

import "gorm.io/gorm"

type Repository struct {
	Wallet      WalletRepository
	Transaction TransactionRepository
}

func New(db *gorm.DB) (Repository, error) {
	return Repository{
		Wallet:      NewAccountRepository(db),
		Transaction: NewTransactionRepository(db),
	}, nil
}
