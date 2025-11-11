package license

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
)

var _ usecase.LicenseUseCase = (*UseCase)(nil)

type UseCase struct {
	signer service.SignatureService
}

func New(signer service.SignatureService) *UseCase {
	return &UseCase{signer: signer}
}
