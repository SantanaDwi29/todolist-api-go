package routes

import (
	"todolist-api/controllers"
	"todolist-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	// Protected routes
	apiRoutes := r.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware())
	{
		// Categories
		apiRoutes.GET("/categories", controllers.GetCategories)
		apiRoutes.POST("/categories", controllers.CreateCategory)
		apiRoutes.DELETE("/categories/:id", controllers.DeleteCategory)

		// Todos
		apiRoutes.GET("/todos", controllers.GetTodos)
		apiRoutes.POST("/todos", controllers.CreateTodo)
		apiRoutes.PUT("/todos/:id", controllers.UpdateTodo)
		apiRoutes.PATCH("/todos/:id/status", controllers.ToggleTodoStatus)
		apiRoutes.DELETE("/todos/:id", controllers.DeleteTodo)
	}
}
