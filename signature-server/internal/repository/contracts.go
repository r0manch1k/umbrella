package repository

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

// LicenseRepository описывает методы для работы с хранением лицензий в базе данных.
type LicenseRepository interface {
	// Save сохраняет лицензию в базу данных.
	// Принимает указатель на entity.License для избежания копирования больших структур.
	Save(ctx context.Context, license *entity.License) error

	// GetByFingerprint возвращает лицензию по её идентификатору и отпечатку устройства.
	// Возвращает указатель на entity.License или nil, если лицензия не найдена.
	GetByFingerprint(ctx context.Context, license, fingerprint string) (*entity.License, error)
}
