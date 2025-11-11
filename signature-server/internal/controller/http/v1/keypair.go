package v1

import (
	httputil "github.com/r0manch1k/umbrella/signature-server/pkg/servers/httpserver"
	"github.com/valyala/fasthttp"
)

func (controller *V1) getPublicKey(ctx *fasthttp.RequestCtx) {
	publicKey, err := controller.keypairUC.GetPublicKey()
	if err != nil {
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, err)

		return
	}

	ctx.SetContentType("application/x-pem-file")
	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.Write(publicKey)
}
