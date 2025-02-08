package http

import (
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/service"
	"github.com/krisnadwipayana07/restful-fintech/pkg/dto"
	"github.com/krisnadwipayana07/restful-fintech/pkg/headers"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Withdraw(c echo.Context) error {
	walletID, err := headers.GetWalletId(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	idempotencyKey, err := headers.GetIdempotencyKey(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	transaction := dto.AmountRequest{}
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	resp, err := h.service.Withdraw(c.Request().Context(), idempotencyKey, walletID, transaction)
	if err != nil {
		return c.JSON(500, dto.BaseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, resp)
}

func (h *TransactionHandler) Deposit(c echo.Context) error {
	walletID, err := headers.GetWalletId(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	idempotencyKey, err := headers.GetIdempotencyKey(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	transaction := dto.AmountRequest{}
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	resp, err := h.service.Deposit(c.Request().Context(), idempotencyKey, walletID, transaction)
	if err != nil {
		return c.JSON(500, dto.BaseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, resp)
}

func (h *TransactionHandler) Transfer(c echo.Context) error {
	walletID, err := headers.GetWalletId(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	idempotencyKey, err := headers.GetIdempotencyKey(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	transaction := dto.TransferRequest{}
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	resp, err := h.service.Transfer(c.Request().Context(), idempotencyKey, walletID, transaction)
	if err != nil {
		return c.JSON(500, dto.BaseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, resp)
}
