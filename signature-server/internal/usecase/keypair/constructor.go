package keypair

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
)

var _ usecase.KeyPairUseCase = (*UseCase)(nil)

type UseCase struct {
	service service.KeyPairService
}

func New(service service.KeyPairService) *UseCase {
	return &UseCase{service: service}
}
