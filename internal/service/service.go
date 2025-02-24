package service

import (
	"context"
	"time"

	"github.com/NeManMarty1/todo-list/internal/logger"
	"github.com/NeManMarty1/todo-list/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
)

type Repository interface {
	CreateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateTask(ctx context.Context, task *models.Task) error
	GetTasks(ctx context.Context, userID int) ([]models.Task, error)
	UpdateTask(ctx context.Context, id, userID int, completed bool) (bool, error)
	DeleteTask(ctx context.Context, id, userID int) (bool, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, password, jwtSecret string) (string, error) {
	existingUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if existingUser != nil {
		logger.Log.WithFields(logrus.Fields{"username": username}).Warn("Пользователь уже существует")
		return "", nil
	}

	user, err := s.repo.CreateUser(ctx, username, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка генерации токена")
		return "", err
	}

	logger.Log.WithFields(logrus.Fields{"user_id": user.ID}).Info("Пользователь зарегистрирован")
	return tokenString, nil
}

func (s *Service) Login(ctx context.Context, username, password, jwtSecret string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil || user.Password != password {
		logger.Log.WithFields(logrus.Fields{"username": username}).Warn("Неверные учетные данные")
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка генерации токена")
		return "", err
	}

	logger.Log.WithFields(logrus.Fields{"user_id": user.ID}).Info("Успешный вход")
	return tokenString, nil
}

func (s *Service) CreateTask(ctx context.Context, task *models.Task) error {
	err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return err
	}
	logger.Log.WithFields(logrus.Fields{"task_id": task.ID, "user_id": task.UserID}).Info("Задача создана")
	return nil
}

func (s *Service) GetTasks(ctx context.Context, userID int) ([]models.Task, error) {
	tasks, err := s.repo.GetTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	logger.Log.WithFields(logrus.Fields{"user_id": userID}).Info("Список задач получен")
	return tasks, nil
}

func (s *Service) UpdateTask(ctx context.Context, id, userID int, completed bool) (bool, error) {
	updated, err := s.repo.UpdateTask(ctx, id, userID, completed)
	if err != nil {
		return false, err
	}
	if !updated {
		logger.Log.WithFields(logrus.Fields{"task_id": id, "user_id": userID}).Warn("Задача не найдена")
		return false, nil
	}
	logger.Log.WithFields(logrus.Fields{"task_id": id, "user_id": userID}).Info("Задача обновлена")
	return true, nil
}

func (s *Service) DeleteTask(ctx context.Context, id, userID int) (bool, error) {
	deleted, err := s.repo.DeleteTask(ctx, id, userID)
	if err != nil {
		return false, err
	}
	if !deleted {
		logger.Log.WithFields(logrus.Fields{"task_id": id, "user_id": userID}).Warn("Задача не найдена")
		return false, nil
	}
	logger.Log.WithFields(logrus.Fields{"task_id": id, "user_id": userID}).Info("Задача удалена")
	return true, nil
}
