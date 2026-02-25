package handler

import (
	"github.com/WardJune/taskflow/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler    *UserHandler
	projectHandler *ProjectHandler
}

func NewRouter(userHandler *UserHandler, projectHandler *ProjectHandler) *Router {
	return &Router{userHandler, projectHandler}
}

func (r *Router) Setup(engine *gin.Engine, jwtSecret string) {
	engine.GET("/health", r.health)

	auth := engine.Group("/api/auth")
	{
		auth.POST("/register", r.userHandler.Register)
		auth.POST("/login", r.userHandler.Login)
	}

	api := engine.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/me", r.userHandler.Me)

		//projects
		projects := api.Group("/projects")
		{
			projects.POST("", r.projectHandler.Create)
			projects.GET("", r.projectHandler.GetMyProjects)
			projects.GET("/:id", r.projectHandler.GetByID)
			projects.POST("/:id/members", r.projectHandler.AddMember)
			projects.POST("/:id/tasks", r.projectHandler.CreateTask)
		}

		//tasks
		tasks := api.Group("/tasks")
		{
			tasks.PATCH("/:task_id", r.projectHandler.UpdateTask)
			tasks.DELETE("/:task_id", r.projectHandler.DeleteTask)
		}
	}
}

func (r *Router) health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "TaskFlow API is running",
	})
}
