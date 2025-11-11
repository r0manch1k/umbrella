package license

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
)

func (r *Repository) DeleteExpired(ctx context.Context, now time.Time) error {
	query, args, err := r.builder.
		Delete("licenses").
		Where(squirrel.Lt{"expires_at": now}).
		ToSql()
	if err != nil {
		return err
	}

	_, execErr := r.db.Pool.Exec(ctx, query, args...)

	return execErr
}
