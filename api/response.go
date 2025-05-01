package api

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	errorMsg string
}

func respondWithError(c echo.Context, status int, msg string) {
	err := ErrorResponse{
		errorMsg: msg,
	}
	c.JSON(status, err)
}