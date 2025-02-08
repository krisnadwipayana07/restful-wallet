package headers

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetWalletId(c echo.Context) (int64, error) {
	walletIdString := c.Request().Header.Get("X-Wallet-ID")
	if walletIdString == "" {
		return 0, errors.New("wallet id not found")
	}

	walletID, err := strconv.ParseInt(walletIdString, 10, 64)
	if err != nil {
		return 0, err
	}

	return walletID, nil
}

func GetIdempotencyKey(c echo.Context) (string, error) {
	// TODO: Format idempotency key
	idempotencyKey := c.Request().Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		return "", errors.New("idempotency key not found")
	}

	return idempotencyKey, nil
}
