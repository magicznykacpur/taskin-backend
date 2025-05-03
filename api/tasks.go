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
	DueUntil    string `json:"due_until"`
}

type TaskRes struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int64  `json:"priority"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DueUntil    string `json:"due_until"`
	UserID      string `json:"user_id"`
}

func mapTaskToTaskRes(task database.Task) TaskRes {
	return TaskRes{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt.Format(time.RFC822),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC822),
		DueUntil:    task.DueUntil.Format(time.RFC822),
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Category:    task.Category,
		UserID:      task.UserID,
	}
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
		createTaskReq.Category == "" ||
		createTaskReq.DueUntil == "" {

		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	dueUntil, err := time.Parse(time.RFC822, createTaskReq.DueUntil)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldnt parse date: %v", err))
	}

	task, err := cfg.DB.CreateTask(
		req.Context(),
		database.CreateTaskParams{
			ID:          uuid.NewString(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DueUntil:    dueUntil,
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

	return c.JSON(http.StatusCreated, mapTaskToTaskRes(task))
}

func (cfg *ApiConfig) HandleGetTaskByID(c echo.Context) error {
	id := c.Param("id")

	task, err := cfg.DB.GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "task not found")
	}

	return c.JSON(http.StatusCreated, mapTaskToTaskRes(task))
}

func (cfg *ApiConfig) HandleGetAllUsersTasks(c echo.Context) error {
	userID := c.Request().Header.Get("userID")

	tasks, err := cfg.DB.GetAllUsersTasks(c.Request().Context(), userID)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "no tasks were found for user with this id")
	}

	tasksRes := []TaskRes{}
	for _, task := range tasks {
		tasksRes = append(tasksRes, mapTaskToTaskRes(task))
	}

	return c.JSON(http.StatusOK, tasksRes)
}

func (cfg *ApiConfig) HandleGetTasksWhereTitleLike(c echo.Context) error {
	title := c.QueryParam("title")
	description := c.QueryParam("description")

	if title != "" && description != "" {
		tasks, err := cfg.DB.GetTaskByTitleAndDescription(
			c.Request().Context(),
			database.GetTaskByTitleAndDescriptionParams{
				Title:       fmt.Sprintf("%%%s%%", title),
				Description: fmt.Sprintf("%%%s%%", description),
			},
		)
		if err != nil {
			return respondWithError(c, http.StatusNotFound, "no tasks found with this title and description")
		}

		tasksRes := []TaskRes{}
		for _, task := range tasks {
			tasksRes = append(tasksRes, mapTaskToTaskRes(task))
		}

		return c.JSON(http.StatusOK, tasksRes)
	}

	if title != "" {
		tasks, err := cfg.DB.GetTasksByTitle(c.Request().Context(), fmt.Sprintf("%%%s%%", title))
		if err != nil {
			return respondWithError(c, http.StatusNotFound, "no tasks found with this title")
		}

		tasksRes := []TaskRes{}
		for _, task := range tasks {
			tasksRes = append(tasksRes, mapTaskToTaskRes(task))
		}

		return c.JSON(http.StatusOK, tasksRes)
	}

	if description != "" {
		tasks, err := cfg.DB.GetTasksByDescription(c.Request().Context(), fmt.Sprintf("%%%s%%", description))
		if err != nil {
			return respondWithError(c, http.StatusNotFound, "no tasks found with this description")
		}

		tasksRes := []TaskRes{}
		for _, task := range tasks {
			tasksRes = append(tasksRes, mapTaskToTaskRes(task))
		}

		return c.JSON(http.StatusOK, tasksRes)
	}

	return respondWithError(c, http.StatusBadRequest, "title or description need to be specified as query parameters")
}

type UpdateTaskReq struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int64  `json:"priority,omitempty"`
	Category    string `json:"category,omitempty"`
	DueUntil    string `json:"due_until,omitempty"`
}

func (cfg *ApiConfig) HandleUpdateTask(c echo.Context) error {
	id := c.Param("id")

	req := c.Request()
	defer req.Body.Close()

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, "couldnt read req bytes")
	}

	var updateTaskReq UpdateTaskReq
	if err := json.Unmarshal(reqBytes, &updateTaskReq); err != nil {
		return respondWithError(c, http.StatusBadRequest, "request body invalid")
	}

	task, err := cfg.DB.GetTaskByID(req.Context(), id)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, "task not found")
	}

	title, description, priority, category, dueUntil, err := retrieveValuesFromTaskUpdateReq(updateTaskReq, task)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("couldnt parse time: %v", err))
	}

	updatedTask, err := cfg.DB.UpdateTaskByID(
		req.Context(),
		database.UpdateTaskByIDParams{
			Title:       title,
			Description: description,
			Priority:    priority,
			Category:    category,
			UpdatedAt:   time.Now(),
			DueUntil:    dueUntil,
			ID:          id,
		},
	)
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("coudlnt update task: %v", err))
	}

	return c.JSON(http.StatusOK, mapTaskToTaskRes(updatedTask))
}

func retrieveValuesFromTaskUpdateReq(updateTaskReq UpdateTaskReq, task database.Task) (string, string, int64, string, time.Time, error) {
	title := updateTaskReq.Title
	if title == "" {
		title = task.Title
	}

	description := updateTaskReq.Description
	if description == "" {
		description = task.Description
	}

	priority := updateTaskReq.Priority
	if priority < 0 {
		priority = task.Priority
	}

	category := updateTaskReq.Category
	if category == "" {
		category = task.Category
	}

	dueUntil, err := time.Parse(time.RFC822, updateTaskReq.DueUntil)
	if err != nil {
		return "", "", -1, "", time.Time{}, err
	}

	return title, description, priority, category, dueUntil, nil
}

type DeleteTaskRes struct {
	Message string `json:"message"`
}

func (cfg *ApiConfig) HandleDeleteTask(c echo.Context) error {
	id := c.Param("id")
	userID := c.Request().Header.Get("userID")

	_, err := cfg.DB.GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return respondWithError(c, http.StatusNotFound, fmt.Sprintf("coudlnt delete task: %v", err))
	}

	err = cfg.DB.DeleteTaskByID(c.Request().Context(), database.DeleteTaskByIDParams{ID: id, UserID: userID})
	if err != nil {
		return respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("coudlnt delete task: %v", err))
	}

	return c.JSON(http.StatusOK, DeleteTaskRes{Message: fmt.Sprintf("task %s deleted successfully", id)})
}
