package main

import (
	"log"
	"os"
	"todo-app/backend/internal/database"
	"todo-app/backend/internal/handlers"
	"todo-app/backend/internal/middleware"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.GinAuthMiddleware())
		{
			protected.GET("/me", authHandler.GetCurrentUser)
			protected.GET("/todos", todoHandler.GetTodos)
			protected.POST("/todos", todoHandler.CreateTodo)
			protected.GET("/todos/:id", todoHandler.GetTodo)
			protected.PUT("/todos/:id", todoHandler.UpdateTodo)
			protected.DELETE("/todos/:id", todoHandler.DeleteTodo)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.GinAuthMiddleware())
		admin.Use(middleware.GinAdminMiddleware())
		{
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
			admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
			admin.GET("/users/:id/todos", adminHandler.GetUserTodos)
		}
	}

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Printf("Default admin credentials - Email: admin@example.com, Password: admin123")
	r.Run(":" + port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
