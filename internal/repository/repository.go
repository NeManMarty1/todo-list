package repository

import (
	"context"
	"database/sql"

	"github.com/NeManMarty1/todo-list/internal/logger"
	"github.com/NeManMarty1/todo-list/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	user := &models.User{Username: username, Password: password}
	err := r.db.QueryRowxContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		username, password).Scan(&user.ID)
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка создания пользователя")
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username=$1", username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка получения пользователя")
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateTask(ctx context.Context, task *models.Task) error {
	err := r.db.QueryRowxContext(ctx, "INSERT INTO tasks (title, completed, user_id) VALUES ($1, $2, $3) RETURNING id",
		task.Title, task.Completed, task.UserID).Scan(&task.ID)
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка создания задачи")
		return err
	}
	return nil
}

func (r *Repository) GetTasks(ctx context.Context, userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.SelectContext(ctx, &tasks, "SELECT * FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка получения задач")
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(ctx context.Context, id, userID int, completed bool) (bool, error) {
	result, err := r.db.ExecContext(ctx, "UPDATE tasks SET completed=$1 WHERE id=$2 AND user_id=$3", completed, id, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка обновления задачи")
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (r *Repository) DeleteTask(ctx context.Context, id, userID int) (bool, error) {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		logger.Log.WithError(err).Error("Ошибка удаления задачи")
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
