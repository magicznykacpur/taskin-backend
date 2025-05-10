package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("coudlnt create user: %v", err))
	}
	if err != nil && strings.Contains(err.Error(), "users.username") {
		return respondWithError(c, http.StatusBadRequest, "user with that username already exists")
	}
	if err != nil && strings.Contains(err.Error(), "users.email") {
		return respondWithError(c, http.StatusBadRequest, "user with that email already exists")
	}

	return c.JSON(http.StatusCreated, CreateUserRes{Username: userReq.Username, Email: userReq.Email})
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	JWTToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *ApiConfig) HandleLoginUser(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()

	requestBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "coudln't read req body bytes")
	}

	var loginReq LoginReq
	if err := json.Unmarshal(requestBytes, &loginReq); err != nil {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	user, err := cfg.DB.GetUserByEmail(c.Request().Context(), loginReq.Email)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "invalid email or password")
	}

	err = auth.ComparePassword(user.HashedPassword, loginReq.Password)
	if err != nil {
		return respondWithError(c, http.StatusUnauthorized, "invalid email or password")
	}

	validRefreshToken, err := cfg.DB.GetValidRefreshTokenForUserId(
		c.Request().Context(),
		database.GetValidRefreshTokenForUserIdParams{
			UserID:    user.ID,
			ExpiresAt: time.Now(),
		},
	)

	if validRefreshToken != (database.RefreshToken{}) {
		jwtToken, err := auth.GenerateJWTToken(user.ID, os.Getenv("JWT_SECRET"), time.Hour)
		if err != nil {
			return respondWithError(c, http.StatusInternalServerError, "couldnt generate jwt token")
		}

		return c.JSON(http.StatusOK, LoginRes{JWTToken: jwtToken, RefreshToken: validRefreshToken.Token})
	} else {
		refreshToken, err := auth.GenerateRefreshToken()
		if err != nil {
			return respondWithError(c, http.StatusInternalServerError, "couldnt generate refresh token")
		}

		err = cfg.DB.CreateRefreshToken(c.Request().Context(), database.CreateRefreshTokenParams{
			UserID:    user.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 31),
		})
		if err != nil {
			return respondWithError(c, http.StatusInternalServerError, "couldnt create refresh token")
		}

		jwtToken, err := auth.GenerateJWTToken(user.ID, os.Getenv("JWT_SECRET"), time.Hour)
		if err != nil {
			return respondWithError(c, http.StatusInternalServerError, "couldnt generate jwt token")
		}

		return c.JSON(http.StatusOK, LoginRes{JWTToken: jwtToken, RefreshToken: refreshToken})
	}
}

type LogoutRes struct {
	Message string `json:"message"`
}

func (cfg *ApiConfig) HandleLogoutUser(c echo.Context) error {
	userID := c.Request().Header.Get("userID")

	user, err := cfg.DB.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "user not found")
	}

	err = cfg.DB.RevokeRefreshToken(c.Request().Context(), user.ID)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "couldnt revoke refresh token")
	}

	return c.JSON(http.StatusOK, LogoutRes{Message: "user successfully logged out"})
}

type UserRes struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (cfg *ApiConfig) HandleGetMe(c echo.Context) error {
	userID := c.Request().Header.Get("userID")
	user, err := cfg.DB.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "user not found")
	}

	return c.JSON(
		200,
		UserRes{
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	)
}

type UpdateUserReq struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (cfg *ApiConfig) HandleUpdateUser(c echo.Context) error {
	userID := c.Request().Header.Get("userID")
	req := c.Request()
	defer req.Body.Close()

	requestBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "coudln't read req body bytes")
	}

	var updateUserReq UpdateUserReq
	if err := json.Unmarshal(requestBytes, &updateUserReq); err != nil {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	if updateUserReq == (UpdateUserReq{}) {
		return respondWithError(c, http.StatusBadRequest, "request body invalid, must provide at least one param to update")
	}

	user, err := cfg.DB.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "user not found")
	}

	email, username, hashedPassword, err := retrieveValuesFromUserUpdateReq(updateUserReq, user)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldnt hash password: %v", err))
	}

	updatedUser, err := cfg.DB.UpdateUserByID(
		req.Context(),
		database.UpdateUserByIDParams{
			ID:             user.ID,
			UpdatedAt:      time.Now(),
			Email:          email,
			Username:       username,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldnt update user: %v", err))
	}

	return c.JSON(
		http.StatusOK,
		UserRes{
			Email:     updatedUser.Email,
			Username:  updatedUser.Username,
			CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
			UpdatedAt: updatedUser.UpdatedAt.Format(time.RFC3339),
		},
	)
}

func retrieveValuesFromUserUpdateReq(updateUserReq UpdateUserReq, user database.User) (string, string, string, error) {
	email := updateUserReq.Email
	if email == "" {
		email = user.Email
	}

	username := updateUserReq.Username
	if username == "" {
		username = user.Username
	}

	password := updateUserReq.Password
	hashedPassword := ""

	if password == "" {
		hashedPassword = user.HashedPassword
	} else {
		hashedPassword, err := auth.HashPassword(password)

		if err != nil {
			return "", "", "", err
		}

		return email, username, string(hashedPassword), nil
	}

	return email, username, string(hashedPassword), nil
}
