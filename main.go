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
	db, err := sql.Open("sqlite", "taskin.db")
	if err != nil {
		log.Fatalf("couldnt open database: %v", err)
	}

	cfg := api.ApiConfig{Port: ":42069", DB: database.New(db)}

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ping")
	})

	e.POST("/api/users", cfg.HandleCreateUser)
	e.POST("/api/login", cfg.HandleLoginUser)

	e.Logger.Fatal(e.Start(cfg.Port))
}
