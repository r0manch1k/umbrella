package di

import (
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/pkg/postgres"
)

type Clients struct {
	PgSQL *postgres.Postgres
}

func initClients(cfg *config.Config) (Clients, error) {
	pg, err := postgres.New(cfg.DB.URL(), false, postgres.MaxPoolSize(cfg.DB.MaxPoolSize))
	if err != nil {
		return Clients{}, err
	}

	return Clients{
		PgSQL: pg,
	}, nil
}
