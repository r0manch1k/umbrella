package v1

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
	httputil "github.com/r0manch1k/umbrella/signature-server/pkg/servers/httpserver"
	"github.com/valyala/fasthttp"
)

func (controller *V1) licenseIssue(ctx *fasthttp.RequestCtx) {
	var req dto.LicenseIssueRequest
	if err := httputil.DecodeJSON(ctx, &req); err != nil {
		httputil.RespondError(ctx, fasthttp.StatusBadRequest, err)

		return
	}

	resp, err := controller.licenseUc.Issue(context.Background(), req)
	if err != nil {
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, err)

		return
	}

	httputil.RespondJSON(ctx, resp)
}

func (controller *V1) licenseVerify(ctx *fasthttp.RequestCtx) {
	var req dto.LicenseVerifyRequest
	if err := httputil.DecodeJSON(ctx, &req); err != nil {
		httputil.RespondError(ctx, fasthttp.StatusBadRequest, err)

		return
	}

	resp, err := controller.licenseUc.Verify(context.Background(), req)
	if err != nil {
		httputil.RespondError(ctx, fasthttp.StatusInternalServerError, err)

		return
	}

	httputil.RespondJSON(ctx, resp)
}
