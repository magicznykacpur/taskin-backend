package api

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func respondWithError(c echo.Context, status int, msg string) error {
	err := ErrorResponse{
		ErrorMessage: msg,
	}
	return c.JSON(status, err)
}
