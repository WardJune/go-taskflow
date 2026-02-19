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
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	//router

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TaskFlow API is running",
		})
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/me", userHandler.Me)
	}

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
