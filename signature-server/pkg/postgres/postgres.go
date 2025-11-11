package postgres

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

var ErrConnAttemptsExceeded = errors.New("postgres - connAttempts == 0")

// Postgres представляет соединение с PostgreSQL и предоставляет пул соединений.
type Postgres struct {
	maxPoolSize  int32
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
}

// New создаёт новый объект Postgres, настраивает пул соединений и подключается к базе.
func New(url string, sslMode bool, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	if sslMode {
		url += "?sslmode=enable"
	} else {
		url += "?sslmode=disable"
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = pg.maxPoolSize

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, ErrConnAttemptsExceeded
	}

	return pg, nil
}

// Close закрывает пул соединений с базой данных.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
