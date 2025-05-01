package api

import "github.com/magicznykacpur/taskin-backend/internal/database"

type ApiConfig struct {
	Port string
	DB *database.Queries
}