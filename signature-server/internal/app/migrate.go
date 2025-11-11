//go:build migrate

package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
	migrationsPath  = "file://migrations"
)

func init() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("config error: %w", err))
	}

	l := logger.New(cfg.Log.Level)

	dbURL := cfg.DB.URL() + "?sslmode=disable"

	var m *migrate.Migrate
	for attempts := defaultAttempts; attempts > 0; attempts-- {
		m, err = migrate.New(migrationsPath, dbURL)
		if err == nil {
			break
		}

		l.Error("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
	}

	if err != nil {
		l.Fatal("Migrate: postgres connect error: %s", err)
	}

	// Выполняем миграции
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Fatal("Migrate: up error: %s", err)
	}

	// Закрываем соединение с миграциями
	if closeErr, _ := m.Close(); closeErr != nil {
		l.Fatal("Migrate: close connection error: %s", closeErr)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("Migrate: no change")

		return
	}

	l.Info("Migrate: up success")
}
