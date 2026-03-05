package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/WardJune/taskflow/internal/handler"
	"github.com/WardJune/taskflow/internal/middleware"
	"github.com/WardJune/taskflow/internal/repository"
	"github.com/WardJune/taskflow/internal/service"
	ws "github.com/WardJune/taskflow/internal/websocket"
	"github.com/WardJune/taskflow/pkg/config"
	"github.com/WardJune/taskflow/pkg/database"
	"github.com/gin-gonic/gin"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	done := make(chan bool, 1)

	cfg := config.Load()

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	hub := ws.NewHub()
	go hub.Run()

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	projectService := service.NewProjectService(projectRepo, taskRepo, userRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo, hub)

	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService, taskService)
	wsHandler := handler.NewWSHandler(hub, projectRepo)

	//router
	engine := gin.Default()
	engine.Use(middleware.CORSmiddleware())
	router := handler.NewRouter(userHandler, projectHandler, wsHandler)
	router.Setup(engine, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: engine,
	}

	go gracefulShutdown(srv, done)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	<-done
	log.Println("Graceful shutodown complete")
}
