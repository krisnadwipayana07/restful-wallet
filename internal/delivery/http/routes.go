package http

import (
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/service"
	"github.com/labstack/echo/v4"
)

const (
	// Ping
	PingPath = "/v1/ping"

	// Transaction
	WithdrawPath = "/v1/withdraw"
	DepositPath  = "/v1/deposit"
	TransferPath = "/v1/transfer"

	// Wallet
	WalletHistoryPath = "/v1/wallet/history"
	WalletBalancePath = "/v1/wallet/balance"
)

func InitHandler(e *echo.Echo, service service.Service) {
	ph := NewPingHandler()
	e.GET(PingPath, ph.Ping)

	th := NewTransactionHandler(service.Transaction)
	e.POST(WithdrawPath, th.Withdraw)
	e.POST(DepositPath, th.Deposit)
	e.POST(TransferPath, th.Transfer)

	wh := NewWalletHandler(service.Wallet)
	e.GET(WalletHistoryPath, wh.WalletHistory)
	e.GET(WalletBalancePath, wh.WalletBalance)
}
