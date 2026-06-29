package routes

import (
	"todolist-api/internal/database"
	"todolist-api/internal/handler"
	"todolist-api/internal/middleware"
	"todolist-api/internal/repository"
	"todolist-api/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	db := database.DB

	// Init Repositories
	authRepo := repository.NewAuthRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	focusRepo := repository.NewFocusRepository(db)
	milestoneRepo := repository.NewMilestoneRepository(db)

	// Init Services
	authService := service.NewAuthService(authRepo)
	todoService := service.NewTodoService(todoRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	focusService := service.NewFocusService(focusRepo)
	milestoneService := service.NewMilestoneService(milestoneRepo)
	analyticsService := service.NewAnalyticsService(todoRepo)

	// Init Handlers
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	focusHandler := handler.NewFocusHandler(focusService)
	milestoneHandler := handler.NewMilestoneHandler(milestoneService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Public routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Protected routes
	apiRoutes := r.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware())
	{
		// Categories
		apiRoutes.GET("/categories", categoryHandler.GetCategories)
		apiRoutes.POST("/categories", categoryHandler.CreateCategory)
		apiRoutes.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Todos
		apiRoutes.GET("/todos", todoHandler.GetTodos)
		apiRoutes.POST("/todos", todoHandler.CreateTodo)
		apiRoutes.PUT("/todos/:id", todoHandler.UpdateTodo)
		apiRoutes.PATCH("/todos/:id/status", todoHandler.ToggleTodoStatus)
		apiRoutes.DELETE("/todos/:id", todoHandler.DeleteTodo)

		// Analytics
		apiRoutes.GET("/analytics", analyticsHandler.GetAnalytics)

		// Focus Session
		apiRoutes.GET("/focus/current", focusHandler.GetCurrentFocusSession)
		apiRoutes.POST("/focus/start", focusHandler.StartFocusSession)
		apiRoutes.POST("/focus/pause", focusHandler.PauseFocusSession)
		apiRoutes.POST("/focus/resume", focusHandler.ResumeFocusSession)
		apiRoutes.POST("/focus/stop", focusHandler.StopFocusSession)

		// Milestones
		apiRoutes.GET("/milestones/next", milestoneHandler.GetNextMilestone)
		apiRoutes.POST("/milestones", milestoneHandler.CreateMilestone)
	}
}
