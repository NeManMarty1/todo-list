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
		logger.Log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	database, err := db.InitDB(cfg.GetDSN())
	if err != nil {
		logger.Log.Fatalf("Ошибка подключения к БД: %v", err)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Log.Infof("Запуск сервера на %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Получен сигнал завершения")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatalf("Ошибка при завершении: %v", err)
	}
	logger.Log.Info("Сервер успешно завершил работу")
}
