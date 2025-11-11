package license

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

func (r *Repository) GetByFingerprint(ctx context.Context, fingerprint string) (*entity.License, error) {
	query, args, err := r.builder.
		Select("fingerprint", "product", "issued_at", "expires_at", "nonce", "activated").
		From("licenses").
		Where(squirrel.Eq{
			"fingerprint": fingerprint,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.Pool.QueryRow(ctx, query, args...)

	var lic entity.License
	if err := row.Scan(
		&lic.Fingerprint,
		&lic.Product,
		&lic.IssuedAt,
		&lic.ExpiresAt,
		&lic.Nonce,
		&lic.Activated,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &lic, nil
}
