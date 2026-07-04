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
	projectRepo := repository.NewProjectRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	focusRepo := repository.NewFocusRepository(db)
	milestoneRepo := repository.NewMilestoneRepository(db)

	// Init Services
	authService := service.NewAuthService(authRepo)
	todoService := service.NewTodoService(todoRepo, projectRepo)
	projectService := service.NewProjectService(projectRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	focusService := service.NewFocusService(focusRepo)
	milestoneService := service.NewMilestoneService(milestoneRepo)
	analyticsService := service.NewAnalyticsService(todoRepo)

	// Init Handlers
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	projectHandler := handler.NewProjectHandler(projectService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	focusHandler := handler.NewFocusHandler(focusService)
	milestoneHandler := handler.NewMilestoneHandler(milestoneService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// API Root Group - Mewajibkan validasi Client ID dan Client Secret
	apiRoot := r.Group("/api")
	apiRoot.Use(middleware.ClientAuthMiddleware())
	{
		// Public routes (Tanpa JWT, tapi butuh Client ID/Secret)
		authRoutes := apiRoot.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
		}

		// Protected routes (Butuh Client ID/Secret + Token JWT User)
		apiRoutes := apiRoot.Group("")
		apiRoutes.Use(middleware.AuthMiddleware())

	{
		// Categories
		apiRoutes.GET("/categories", categoryHandler.GetCategories)
		apiRoutes.POST("/categories", categoryHandler.CreateCategory)
		apiRoutes.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Projects
		apiRoutes.GET("/projects", projectHandler.GetAllProjects)
		apiRoutes.POST("/projects", projectHandler.CreateProject)
		apiRoutes.GET("/projects/:id", projectHandler.GetProjectByID)
		apiRoutes.DELETE("/projects/:id", projectHandler.DeleteProject)

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
}

