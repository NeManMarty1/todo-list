package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NeManMarty1/todo-list/internal/config"
	"github.com/NeManMarty1/todo-list/internal/db"
	"github.com/NeManMarty1/todo-list/internal/handlers"
	"github.com/NeManMarty1/todo-list/internal/logger"
	"github.com/NeManMarty1/todo-list/internal/middleware"
	"github.com/NeManMarty1/todo-list/internal/repository"
	"github.com/NeManMarty1/todo-list/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Log.WithError(err).Fatal("failed to load config")
	}

	database, err := db.InitDB(cfg.GetDSN())
	if err != nil {
		logger.Log.WithError(err).Fatal("failed to connect to database")
	}
	defer database.Close()

	repo := repository.NewRepository(database)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc, cfg)

	r := gin.Default()

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	authorized := r.Group("/").Use(middleware.AuthMiddleware(cfg))
	{
		authorized.GET("/profile", h.Profile)
		authorized.POST("/tasks", h.CreateTask)
		authorized.GET("/tasks", h.GetTasks)
		authorized.PUT("/tasks/:id", h.UpdateTask)
		authorized.DELETE("/tasks/:id", h.DeleteTask)
	}

	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	ctx := context.Background()

	go func() {
		logger.Log.WithField("port", cfg.Server.Port).Info("starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.WithError(err).Fatal("server failed to start")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logger.Log.Info("shutting down server...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.WithError(err).Error("server shutdown failed")
		return
	}

	logger.Log.Info("server gracefully stopped")
}
