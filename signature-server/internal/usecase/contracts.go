package usecase

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

type LicenseUseCase interface {
	// Issue — метод для выдачи лицензий.
	Issue(ctx context.Context, req dto.LicenseIssueRequest) (dto.LicenseIssueResponse, error)

	// Verify — метод для проверки подлинности лицензий.
	Verify(ctx context.Context, req dto.LicenseVerifyRequest) (dto.LicenseVerifyResponse, error)
}

type KeyPairUseCase interface {
	GenerateAndSaveKeyPair() error
	GetPublicKey() ([]byte, error)
}
