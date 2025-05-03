package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/magicznykacpur/taskin-backend/internal/database"
)

type CreateTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int64  `json:"priority"`
	Category    string `json:"category"`
}

type TaskRes struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int64  `json:"priority"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	UserID      string `json:"user_id"`
}

func (cfg *ApiConfig) HandleCreateTask(c echo.Context) error {
	req := c.Request()
	defer req.Body.Close()

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "couldnt read req bytes")
	}

	var createTaskReq CreateTaskReq
	if err := json.Unmarshal(reqBytes, &createTaskReq); err != nil {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	if createTaskReq.Title == "" ||
		createTaskReq.Description == "" ||
		createTaskReq.Priority < 0 ||
		createTaskReq.Category == "" {

		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	task, err := cfg.DB.CreateTask(
		req.Context(),
		database.CreateTaskParams{
			ID:          uuid.NewString(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       createTaskReq.Title,
			Description: createTaskReq.Description,
			Priority:    createTaskReq.Priority,
			Category:    createTaskReq.Category,
			UserID:      req.Header.Get("userID"),
		},
	)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldnt create task: %v", err))
	}

	return c.JSON(
		201,
		TaskRes{
			ID:          task.ID,
			CreatedAt:   task.CreatedAt.Format(time.UnixDate),
			UpdatedAt:   task.UpdatedAt.Format(time.UnixDate),
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Category:    task.Category,
			UserID:      task.UserID,
		})
}

func (cfg *ApiConfig) HandleGetTaskByID(c echo.Context) error {
	id := c.Param("id")

	task, err := cfg.DB.GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "task not found")
	}

	return c.JSON(
		200,
		TaskRes{
			ID:          task.ID,
			CreatedAt:   task.CreatedAt.Format(time.UnixDate),
			UpdatedAt:   task.UpdatedAt.Format(time.UnixDate),
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Category:    task.Category,
			UserID:      task.UserID,
		},
	)
}

func (cfg *ApiConfig) HandleGetAllUsersTasks(c echo.Context) error {
	userID := c.Request().Header.Get("userID")

	tasks, err := cfg.DB.GetAllUsersTasks(c.Request().Context(), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "no tasks were found for user with this id")
	}

	tasksRes := []TaskRes{}
	for _, task := range tasks {
		tasksRes = append(tasksRes, TaskRes{
			ID:          task.ID,
			CreatedAt:   task.CreatedAt.Format(time.UnixDate),
			UpdatedAt:   task.UpdatedAt.Format(time.UnixDate),
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Category:    task.Category,
			UserID:      task.UserID,
		})
	}

	return c.JSON(200, tasksRes)
}
