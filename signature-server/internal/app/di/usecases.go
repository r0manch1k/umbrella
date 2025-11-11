package di

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase/license"
)

type UseCases struct {
	License usecase.LicenseUseCase
}

func initUseCases(services Services) UseCases {
	return UseCases{
		License: license.New(services.Signature),
	}
}
