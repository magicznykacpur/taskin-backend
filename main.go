package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/magicznykacpur/taskin-backend/api"
	"github.com/magicznykacpur/taskin-backend/internal/database"
	_ "modernc.org/sqlite"
)

func loadEnvVars() {
	data, err := os.ReadFile(".env")
	if err != nil {
		log.Fatalf("coudlnt open env file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")
		os.Setenv(parts[0], parts[1])
	}
}

func main() {
	loadEnvVars()

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

	e.POST("/api/signup", cfg.HandleCreateUser)
	e.POST("/api/login", cfg.HandleLoginUser)
	e.POST("/api/logout", cfg.HandleLogoutUser, cfg.LoggedInMiddleware)
	e.GET("/api/me", cfg.HandleGetMe, cfg.LoggedInMiddleware)
	e.PUT("/api/users", cfg.HandleUpdateUser, cfg.LoggedInMiddleware)

	e.POST("/api/tasks", cfg.HandleCreateTask, cfg.LoggedInMiddleware)
	e.GET("/api/tasks", cfg.HandleGetAllUsersTasks, cfg.LoggedInMiddleware)
	e.GET("/api/tasks/:id", cfg.HandleGetTaskByID, cfg.LoggedInMiddleware)
	e.PUT("/api/tasks/:id", cfg.HandleUpdateTask, cfg.LoggedInMiddleware)
	e.DELETE("/api/tasks/:id", cfg.HandleDeleteTask, cfg.LoggedInMiddleware)
	e.GET("/api/tasks/search", cfg.HandleGetTasksWhereTitleLike, cfg.LoggedInMiddleware)

	e.Logger.Fatal(e.Start(cfg.Port))
}
