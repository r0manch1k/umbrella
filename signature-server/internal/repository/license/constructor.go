package license

import (
	"github.com/Masterminds/squirrel"
	"github.com/r0manch1k/umbrella/signature-server/pkg/postgres"
)

type Repository struct {
	db      *postgres.Postgres
	builder squirrel.StatementBuilderType
}

func New(db *postgres.Postgres) *Repository {
	return &Repository{
		db:      db,
		builder: db.Builder,
	}
}
