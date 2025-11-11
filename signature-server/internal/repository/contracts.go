package repository

import (
	"context"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

type LicenseRepository interface {
	Save(ctx context.Context, license entity.License) error
	GetByUserAndFingerprint(ctx context.Context, userID, hwFingerprint string) (*entity.License, error)
	DeleteExpired(ctx context.Context, now time.Time) error
}
