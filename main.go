package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magicznykacpur/taskin-backend/api"
	"github.com/magicznykacpur/taskin-backend/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "taskin")
	if err != nil {
		log.Fatalf("couldnt open database: %v", err)
	}

	queries := database.New(db)
	cfg := api.ApiConfig{Port: ":42069", DB: queries}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ping")
	})

	e.POST("/api/users", cfg.HandleCreateUser)

	e.Logger.Fatal(e.Start(cfg.Port))
}
