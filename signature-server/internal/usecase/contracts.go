package usecase

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

// LicenseUseCase описывает бизнес-логику работы с лицензиями.
type LicenseUseCase interface {
	// Issue выдаёт лицензию на основе запроса dto.LicenseIssueRequest.
	// Возвращает dto.LicenseIssueResponse с данными лицензии или ошибку.
	Issue(ctx context.Context, req dto.LicenseIssueRequest) (dto.LicenseIssueResponse, error)

	// Verify проверяет подлинность лицензии на основе запроса dto.LicenseVerifyRequest.
	// Возвращает dto.LicenseVerifyResponse с результатом проверки или ошибку.
	Verify(ctx context.Context, req dto.LicenseVerifyRequest) (dto.LicenseVerifyResponse, error)

	// GetPublicKey возвращает публичный ключ сервиса для проверки лицензий.
	GetPublicKey() ([]byte, error)
}
