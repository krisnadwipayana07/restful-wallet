package model

import "time"

type Wallet struct {
	ID             int64      `gorm:"column:id"`
	Name           string     `gorm:"column:wallet_name"`
	CurrentBalance string     `gorm:"column:wallet_curr_balance"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (Wallet) TableName() string {
	return "wallet_table"
}
