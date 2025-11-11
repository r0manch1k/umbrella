package license

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

func (r *Repository) Save(ctx context.Context, license *entity.License) error {
	query, args, err := r.builder.
		Insert("licenses").
		Columns("fingerprint", "product", "issued_at", "expires_at", "nonce", "activated").
		Values(license.Fingerprint, license.Product, license.IssuedAt, license.ExpiresAt, license.Nonce, license.Activated).
		Suffix(`ON CONFLICT (fingerprint) DO UPDATE
			SET product = EXCLUDED.product,
			    issued_at = EXCLUDED.issued_at,
			    expires_at = EXCLUDED.expires_at,
			    nonce = EXCLUDED.nonce,
			    activated = EXCLUDED.activated`).
		ToSql()
	if err != nil {
		return err
	}

	_, execErr := r.db.Pool.Exec(ctx, query, args...)

	return execErr
}
