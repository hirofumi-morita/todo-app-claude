package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	DB *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{DB: db}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	userCtx, ok := middleware.GetUserFromGinContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		 FROM todos WHERE user_id = $1 ORDER BY created_at DESC`,
		userCtx.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
			&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan todo"})
			return
		}
		todos = append(todos, todo)
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	userCtx, ok := middleware.GetUserFromGinContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var todo models.Todo
	err = h.DB.QueryRow(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		 FROM todos WHERE id = $1 AND user_id = $2`,
		todoID, userCtx.UserID,
	).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userCtx, ok := middleware.GetUserFromGinContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	var todo models.Todo
	err := h.DB.QueryRow(
		`INSERT INTO todos (user_id, title, description, completed)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, description, completed, created_at, updated_at`,
		userCtx.UserID, req.Title, req.Description, req.Completed,
	).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userCtx, ok := middleware.GetUserFromGinContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var req models.TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var todo models.Todo
	err = h.DB.QueryRow(
		`UPDATE todos
		 SET title = $1, description = $2, completed = $3, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, title, description, completed, created_at, updated_at`,
		req.Title, req.Description, req.Completed, todoID, userCtx.UserID,
	).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userCtx, ok := middleware.GetUserFromGinContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	result, err := h.DB.Exec(
		"DELETE FROM todos WHERE id = $1 AND user_id = $2",
		todoID, userCtx.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
