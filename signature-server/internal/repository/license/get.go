package license

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

func (r *Repository) GetByUserAndFingerprint(ctx context.Context, userID, hwFingerprint string) (*entity.License, error) {
	query, args, err := r.builder.
		Select("user_id", "product", "issued_at", "expires_at", "hw_fingerprint", "nonce").
		From("licenses").
		Where(squirrel.Eq{
			"user_id":        userID,
			"hw_fingerprint": hwFingerprint,
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.Pool.QueryRow(ctx, query, args...)

	var lic entity.License
	if err := row.Scan(
		&lic.UserID,
		&lic.Product,
		&lic.IssuedAt,
		&lic.ExpiresAt,
		&lic.HWFingerprint,
		&lic.Nonce,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &lic, nil
}
