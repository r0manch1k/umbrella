package v1

import (
	"github.com/fasthttp/router"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
)

func NewLicenseRoutes(
	v1Group *router.Group,
	l logger.Interface,
	licenseUC usecase.LicenseUseCase,
) {
	controller := &V1{l: l, licenseUc: licenseUC}

	licenseGroup := v1Group.Group("/license")
	{
		licenseGroup.POST("/issue", controller.licenseIssue)
		licenseGroup.POST("/verify", controller.licenseVerify)
	}
}

func NewKeypairRoutes(
	v1Group *router.Group,
	l logger.Interface,
	keypairUC usecase.KeyPairUseCase,
) {
	controller := &V1{l: l, keypairUC: keypairUC}

	keypairGroup := v1Group.Group("/keypair")
	{
		keypairGroup.GET("/public.pem", controller.getPublicKey)
	}
}
