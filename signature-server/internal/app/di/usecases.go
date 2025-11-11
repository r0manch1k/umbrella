package di

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase/keypair"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase/license"
)

type UseCases struct {
	KeyPair usecase.KeyPairUseCase
	License usecase.LicenseUseCase
}

func initUseCases(services Services) UseCases {
	return UseCases{
		KeyPair: keypair.New(services.KeyPair),
		License: license.New(services.Signature),
	}
}
