package http

import (
	"github.com/fasthttp/router"
	v1 "github.com/r0manch1k/umbrella/signature-server/internal/controller/http/v1"
	"github.com/r0manch1k/umbrella/signature-server/internal/usecase"
	"github.com/r0manch1k/umbrella/signature-server/pkg/logger"
	"github.com/valyala/fasthttp"
)

type Router struct {
	Router *router.Router
}

// NewRouter - инициализирует все маршруты приложения.
func NewRouter(
	l logger.Interface,
	licenseUC usecase.LicenseUseCase,
	keypairUC usecase.KeyPairUseCase,
) *Router {
	rt := &Router{
		Router: router.New(),
	}

	// ============== Utils Routes ==============

	rt.Router.GET("/health", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.SetBodyString("OK!")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
	})

	// ============== V1 Routes ==============

	v1Group := rt.Router.Group("/v1")
	{
		v1.NewLicenseRoutes(v1Group, l, licenseUC)
		v1.NewKeypairRoutes(v1Group, l, keypairUC)
	}

	return rt
}
