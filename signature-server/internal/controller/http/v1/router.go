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
		licenseGroup.POST("/issue", controller.issue)
		licenseGroup.POST("/verify", controller.verify)
		licenseGroup.GET("/public.pem", controller.getPublicKey)
	}
}
