package migrate

import (
	"context"

	"github.com/NeManMarty1/todo-list/internal/config"
	"github.com/NeManMarty1/todo-list/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, cfg *config.Config) {
	logger.Init()

	sourceURL := "file:///migrations" 
	databaseURL := cfg.GetDSN()

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		logger.Log.WithContext(ctx).Fatalf("Ошибка создании миграции: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.WithContext(ctx).Fatalf("Ошибка применения миграций: %v", err)
	}

	logger.Log.WithContext(ctx).Info("Миграции успешно применены")
}
