//go:build migrate

package app

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/r0manch1k/umbrella/signature-server/config"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("config error: %s", err))
	}

	dbUrl := cfg.DB.URL() + "?sslmode=disable"

	for attempts > 0 {
		m, err = migrate.New("file://migrations", dbUrl)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		log.Printf(dbUrl)

		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Fatalf("Migrate: postgres connection error: %s", err)
		}
	}(m)

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")

		return
	}

	log.Printf("Migrate: up success")
}
