package main

import (
	"log"

	"github.com/WardJune/taskflow/internal/handler"
	"github.com/WardJune/taskflow/internal/middleware"
	"github.com/WardJune/taskflow/internal/repository"
	"github.com/WardJune/taskflow/internal/service"
	"github.com/WardJune/taskflow/pkg/config"
	"github.com/WardJune/taskflow/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	projectService := service.NewProjectService(projectRepo, taskRepo, userRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo)

	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService, taskService)

	//router
	engine := gin.New()
	engine.Use(middleware.CORSmiddleware())
	router := handler.NewRouter(userHandler, projectHandler)
	router.Setup(engine, cfg.JWTSecret)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := engine.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
