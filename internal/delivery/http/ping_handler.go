package http

import "github.com/labstack/echo/v4"

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}
