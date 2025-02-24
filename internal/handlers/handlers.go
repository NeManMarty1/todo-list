package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NeManMarty1/todo-list/internal/config"
	"github.com/NeManMarty1/todo-list/internal/logger"
	"github.com/NeManMarty1/todo-list/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Register(ctx context.Context, username, password, jwtSecret string) (string, error)
	Login(ctx context.Context, username, password, jwtSecret string) (string, error)
	CreateTask(ctx context.Context, task *models.Task) error
	GetTasks(ctx context.Context, userID int) ([]models.Task, error)
	UpdateTask(ctx context.Context, id, userID int, completed bool) (bool, error)
	DeleteTask(ctx context.Context, id, userID int) (bool, error)
}

type Handler struct {
	svc Service
	cfg *config.Config
}

func NewHandler(svc Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}

func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Log.WithError(err).Warn("Неверный формат данных в /register")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	token, err := h.svc.Register(ctx, user.Username, user.Password, h.cfg.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if token == "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Log.WithError(err).Warn("Неверный формат данных в /login")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	token, err := h.svc.Login(ctx, user.Username, user.Password, h.cfg.JWT.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные данные"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetInt("userID")
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *Handler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		logger.Log.WithError(err).Warn("Неверный формат данных в /tasks")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}
	task.UserID = c.GetInt("userID")

	if err := h.svc.CreateTask(ctx, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) GetTasks(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt("userID")
	tasks, err := h.svc.GetTasks(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("userID")
	var update struct {
		Completed bool `json:"completed"`
	}
	if err := c.BindJSON(&update); err != nil {
		logger.Log.WithError(err).Warn("Неверный формат данных в /tasks/:id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	updated, err := h.svc.UpdateTask(ctx, id, userID, update.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Задача обновлена"})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("userID")

	deleted, err := h.svc.DeleteTask(ctx, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}
