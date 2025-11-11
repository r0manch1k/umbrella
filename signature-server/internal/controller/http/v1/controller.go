package v1

import (
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
)

type V1 struct {
	l         logger.Interface
	licenseUc usecase.LicenseUseCase
	keypairUC usecase.KeyPairUseCase
}
