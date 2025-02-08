package http

import (
	"github.com/krisnadwipayana07/restful-fintech/internal/domain/service"
	"github.com/krisnadwipayana07/restful-fintech/pkg/dto"
	"github.com/krisnadwipayana07/restful-fintech/pkg/headers"
	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) WalletHistory(c echo.Context) error {
	walletID, err := headers.GetWalletId(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	transactionList, err := h.service.WalletHistory(c.Request().Context(), walletID)
	if err != nil {
		return c.JSON(500, dto.BaseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, transactionList)
}

func (h *WalletHandler) WalletBalance(c echo.Context) error {
	walletID, err := headers.GetWalletId(c)
	if err != nil {
		return c.JSON(400, dto.BaseError{
			Message: err.Error(),
		})
	}

	balance, err := h.service.WalletBalance(c.Request().Context(), walletID)
	if err != nil {
		return c.JSON(500, dto.BaseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, balance)
}
