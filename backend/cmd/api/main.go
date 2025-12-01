package main

import (
	"log"
	"net/http"
	"os"
	"todo-app/backend/internal/database"
	"todo-app/backend/internal/handlers"
	"todo-app/backend/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "todoapp"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := middleware.CreateDefaultAdmin(db); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	authHandler := handlers.NewAuthHandler(db)
	todoHandler := handlers.NewTodoHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")

	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/me", authHandler.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	protected.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	protected.HandleFunc("/todos/{id}", todoHandler.GetTodo).Methods("GET")
	protected.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	protected.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.AdminMiddleware)
	admin.HandleFunc("/users", adminHandler.GetAllUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", adminHandler.GetUser).Methods("GET")
	admin.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/users/{id}/role", adminHandler.UpdateUserRole).Methods("PUT")
	admin.HandleFunc("/users/{id}/todos", adminHandler.GetUserTodos).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Printf("Default admin credentials - Email: admin@example.com, Password: admin123")
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
