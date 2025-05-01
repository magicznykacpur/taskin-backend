package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/magicznykacpur/taskin-backend/auth"
	"github.com/magicznykacpur/taskin-backend/internal/database"
)

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (cfg *ApiConfig) HandleCreateUser(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()

	requestBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "coudln't read req body bytes")
	}

	var userReq CreateUserReq
	if err := json.Unmarshal(requestBytes, &userReq); err != nil {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	if userReq.Username == "" || userReq.Password == "" || userReq.Email == "" {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	hash, err := auth.HashPassword(userReq.Password)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "coudlnt hash password")
	}

	err = cfg.DB.CreateUser(c.Request().Context(), database.CreateUserParams{
		ID:             uuid.NewString(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          userReq.Email,
		Username:       userReq.Username,
		HashedPassword: string(hash),
	})
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldn't create user: %v", err))
	}

	return c.JSON(http.StatusCreated, CreateUserRes{Username: userReq.Username, Email: userReq.Email})
}
