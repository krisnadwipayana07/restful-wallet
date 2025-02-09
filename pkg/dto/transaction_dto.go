package dto

import (
	"time"

	"github.com/krisnadwipayana07/restful-fintech/internal/domain/model"
	"github.com/shopspring/decimal"
)

type AmountRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

type TransferRequest struct {
	ReceiverWalletID int64           `json:"receiver_wallet_id"`
	Amount           decimal.Decimal `json:"amount"`
}

type TransactionResponse struct {
	TransactionID int64 `json:"transaction_id"`
}

type TransactionDetailResponse struct {
	TransactionID int64           `json:"transaction_id"`
	WalletID      int64           `json:"wallet_id"`
	Type          int16           `json:"type"`
	IsDebit       bool            `json:"is_debit"`
	Value         decimal.Decimal `json:"value"`
	Remarks       string          `json:"remarks"`
	CreatedAt     time.Time       `json:"created_at"`
}

func NewTransactionListResponse(transactions []model.Transaction) []TransactionDetailResponse {
	var resp []TransactionDetailResponse
	for _, transaction := range transactions {
		resp = append(resp, TransactionDetailResponse{
			TransactionID: transaction.ID,
			WalletID:      transaction.WalletID,
			Type:          transaction.Type,
			IsDebit:       transaction.IsDebit,
			Value:         transaction.Value,
			Remarks:       transaction.Remarks,
			CreatedAt:     transaction.CreatedAt,
		})
	}
	return resp
}
