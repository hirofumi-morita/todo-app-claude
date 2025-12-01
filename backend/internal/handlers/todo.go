package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/models"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	DB *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{DB: db}
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(
		`SELECT id, user_id, title, description, completed, created_at, updated_at
		 FROM todos WHERE user_id = $1 ORDER BY created_at DESC`,
		userCtx.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
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
			http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
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
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
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
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec(
		"DELETE FROM todos WHERE id = $1 AND user_id = $2",
		todoID, userCtx.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}
