package api

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/magicznykacpur/taskin-backend/auth"
	"github.com/magicznykacpur/taskin-backend/internal/database"
)

func (cfg *ApiConfig) LoggedInMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		refreshToken, bearerToken, err := auth.GetAuthTokensFromHeaders(c.Request().Header)
		if err != nil {
			return respondWithError(c, http.StatusUnauthorized, "missing authorization")
		}

		userID, err := auth.ValidateJWTToken(bearerToken, os.Getenv("JWT_SECRET"))
		if err != nil && err.Error() == "token has invalid claims: token is expired" {
			dbRefreshToken, err := cfg.DB.GetValidRefreshTokenByValue(
				c.Request().Context(),
				database.GetValidRefreshTokenByValueParams{
					ExpiresAt: time.Now(),
					Token:     refreshToken,
				},
			)
			if err != nil {
				return respondWithError(c, http.StatusUnauthorized, "invalid refresh or jwt token")
			}

			freshJWT, err := auth.GenerateJWTToken(dbRefreshToken.UserID, os.Getenv("JWT_SECRET"), time.Hour)
			c.SetCookie(&http.Cookie{Name: "jwt_token", Value: freshJWT, Path: "/"})
			c.Request().Header.Set("userID", dbRefreshToken.UserID)

			return next(c)
		}

		c.Request().Header.Set("userID", userID)
		return next(c)
	}
}
