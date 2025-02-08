package service

import (
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db          *gorm.DB
	Transaction TransactionService
	Wallet      WalletService
}

func New(repo repository.Repository, db *gorm.DB, redis *redis.Client) (Service, error) {
	return Service{
		db:          db,
		Transaction: NewTransactionService(db, redis, repo.Transaction, repo.Wallet),
		Wallet:      NewWalletService(repo.Transaction, repo.Wallet),
	}, nil
}
