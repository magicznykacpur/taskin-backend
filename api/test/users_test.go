package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/magicznykacpur/taskin-backend/api"
	"github.com/magicznykacpur/taskin-backend/internal/database"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

var (
	validUserReq     = `{"username":"test user","password":"password","email":"email@test.com"}`
	validLoginReq    = `{"password":"password","email":"email@test.com"}`
	invalidUserReq   = `{"username":"test user"}`
	malformedUserReq = `{"username`
)

const createUsersTable = `CREATE TABLE users (
    id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL
);`

const createRefreshTokensTable = `CREATE TABLE refresh_tokens (
    user_id TEXT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_revoked INTEGER NOT NULL,
    expires_at TIMESTAMP NOT NULL
);`

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(context.Background(), createUsersTable)
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(context.Background(), createRefreshTokensTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupEcho(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
func TestValidUserReq(t *testing.T) {
	c, rec := setupEcho(http.MethodPost, "/api/users", validUserReq)
	db, err := setupDB()
	if err != nil {
		log.Fatalf("coudlnt create database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}

	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	res := rec.Result()
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var userRes api.CreateUserRes
	if err := json.Unmarshal(resBytes, &userRes); err != nil {
		log.Fatalf("couldnt unmarshall res body: %v", err)
	}

	assert.Equal(t, "test user", userRes.Username)

	users, err := cfg.DB.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("coudlnt retrieve users: %v", err)
	}

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "test user", users[0].Username)
	assert.Equal(t, "email@test.com", users[0].Email)
}

func TestInvalidUserReq(t *testing.T) {
	c, rec := setupEcho(http.MethodPost, "/api/users", invalidUserReq)
	db, err := setupDB()
	if err != nil {
		log.Fatalf("coudlnt create database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}
	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	res := rec.Result()
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var errorRes api.ErrorResponse
	if err := json.Unmarshal(resBytes, &errorRes); err != nil {
		log.Fatalf("couldnt unmarshall res body: %v", err)
	}

	assert.Equal(t, "request body invalid", errorRes.ErrorMessage)
}

func TestMalformedRequestBody(t *testing.T) {
	c, rec := setupEcho(http.MethodPost, "/api/users", malformedUserReq)
	db, err := setupDB()
	if err != nil {
		log.Fatalf("coudlnt create database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}
	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	res := rec.Result()
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var errorRes api.ErrorResponse
	if err := json.Unmarshal(resBytes, &errorRes); err != nil {
		log.Fatalf("couldnt unmarshall res body: %v", err)
	}

	assert.Equal(t, "request body invalid", errorRes.ErrorMessage)
}

func TestUniqueUser(t *testing.T) {
	c, rec := setupEcho(http.MethodPost, "/api/users", validUserReq)
	db, err := setupDB()
	if err != nil {
		log.Fatalf("coudlnt create database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}
	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	res := rec.Result()
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var userRes api.CreateUserRes
	if err := json.Unmarshal(resBytes, &userRes); err != nil {
		log.Fatalf("couldnt unmarshall res body: %v", err)
	}

	assert.Equal(t, "test user", userRes.Username)

	users, err := cfg.DB.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("coudlnt retrieve users: %v", err)
	}

	assert.Equal(t, 1, len(users))
	assert.Equal(t, "test user", users[0].Username)

	c, rec = setupEcho(http.MethodPost, "/api/users", validUserReq)

	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	res = rec.Result()
	defer res.Body.Close()
	resBytes, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var errorRes api.ErrorResponse
	if err := json.Unmarshal(resBytes, &errorRes); err != nil {
		log.Fatalf("couldnt unmarshall res body: %v", err)
	}

	assert.Equal(t, "couldn't create user: constraint failed: UNIQUE constraint failed: users.email (2067)", errorRes.ErrorMessage)
}

func TestLoginUser(t *testing.T) {
	c, _ := setupEcho(http.MethodPost, "/api/login", validUserReq)
	db, err := setupDB()

	if err != nil {
		log.Fatalf("coudlnt create database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}

	err = cfg.HandleCreateUser(c)
	assert.NoError(t, err)

	c, rec := setupEcho(http.MethodPost, "/api/login", validLoginReq)

	err = cfg.HandleLoginUser(c)
	assert.NoError(t, err)

	res := rec.Result()
	defer res.Body.Close()
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("couldnt read res body bytes: %v", err)
	}

	var loginRes api.LoginRes
	if err := json.Unmarshal(resBytes, &loginRes); err != nil {
		log.Fatalf("couldnt unmarshall res bytes: %v", err)
	}

	assert.NotNil(t, loginRes.JWTToken)
	assert.NotNil(t, loginRes.RefreshToken)

	users, err := cfg.DB.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("couldnt retrieve user: %v", err)
	}

	refreshToken, err := cfg.DB.GetValidRefreshTokenForUserId(context.Background(),
		database.GetValidRefreshTokenForUserIdParams{
			UserID:    users[0].ID,
			ExpiresAt: time.Now(),
		},
	)
	if err != nil {
		log.Fatalf("couldnt retrieve refresh token: %v", err)
	}

	assert.NotNil(t, refreshToken)
	assert.Equal(t, users[0].ID, refreshToken.UserID)
}
