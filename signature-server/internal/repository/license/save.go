package license

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

func (r *Repository) Save(ctx context.Context, license entity.License) error {
	query, args, err := r.builder.
		Insert("licenses").
		Columns("user_id", "product", "issued_at", "expires_at", "hw_fingerprint", "nonce").
		Values(license.UserID, license.Product, license.IssuedAt, license.ExpiresAt, license.HWFingerprint, license.Nonce).
		Suffix("ON CONFLICT (user_id, hw_fingerprint) DO UPDATE SET product = EXCLUDED.product, issued_at = EXCLUDED.issued_at, expires_at = EXCLUDED.expires_at, nonce = EXCLUDED.nonce").
		ToSql()
	if err != nil {
		return err
	}

	_, execErr := r.db.Pool.Exec(ctx, query, args...)

	return execErr
}
