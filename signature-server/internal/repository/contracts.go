package repository

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

type LicenseRepository interface {
	Save(ctx context.Context, license *entity.License) error
	GetByFingerprint(ctx context.Context, license, fingerprint string) (*entity.License, error)
}
