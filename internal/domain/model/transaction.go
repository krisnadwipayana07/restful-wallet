package model

import "time"

type Transaction struct {
	ID        int64     `gorm:"column:id"`
	WalletID  int64     `gorm:"column:wallet_id"`
	Type      int16     `gorm:"column:trc_type"`
	IsDebit   bool      `gorm:"column:trc_is_debit"`
	Value     string    `gorm:"column:trc_value"`
	Remarks   string    `gorm:"column:trc_remarks"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Transaction) TableName() string {
	return "transaction_table"
}
