package v1

import (
	"context"
	"errors"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
	httputil "github.com/r0manch1k/umbrella/signature-server/pkg/servers/httpserver"
	"github.com/valyala/fasthttp"
)

func (controller *V1) issue(ctx *fasthttp.RequestCtx) {
	var req dto.LicenseIssueRequest
	if err := httputil.DecodeJSON(ctx, &req); err != nil {
		httputil.RespondError(ctx, fasthttp.StatusBadRequest, err)

		return
	}

	resp, err := controller.licenseUc.Issue(context.Background(), req)
	if err != nil {
		switch {
		case errors.Is(err, exception.ErrLicenseAlreadyActivatedAndNotExpired):
			httputil.RespondError(ctx, fasthttp.StatusConflict, err)

		default:
			controller.l.Error("issue: internal error: %v", err)
			httputil.RespondError(ctx, fasthttp.StatusInternalServerError, exception.ErrInternal)
		}

		return
	}

	httputil.RespondJSON(ctx, resp)
}

func (controller *V1) verify(ctx *fasthttp.RequestCtx) {
	var req dto.LicenseVerifyRequest
	if err := httputil.DecodeJSON(ctx, &req); err != nil {
		httputil.RespondError(ctx, fasthttp.StatusBadRequest, err)

		return
	}

	resp, err := controller.licenseUc.Verify(context.Background(), req)
	if err != nil {
		controller.l.Error("verify: internal error: %v", err)
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, exception.ErrInternal)

		return
	}

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.Response.SetBodyString(resp.Signature)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (controller *V1) getPublicKey(ctx *fasthttp.RequestCtx) {
	publicKey, err := controller.licenseUc.GetPublicKey()
	if err != nil {
		controller.l.Error("getPublicKey: failed to get public key: %v", err)
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, err)

		return
	}

	if _, err := ctx.Write(publicKey); err != nil {
		controller.l.Error("getPublicKey: failed to write response: %v", err)
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, exception.ErrInternal)

		return
	}

	ctx.SetContentType("application/x-pem-file")
	ctx.SetStatusCode(fasthttp.StatusOK)
}
