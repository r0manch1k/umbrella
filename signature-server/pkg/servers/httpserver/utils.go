package httpserver

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func DecodeJSON(ctx *fasthttp.RequestCtx, dst any) error {
	return json.Unmarshal(ctx.PostBody(), dst)
}

func RespondJSON(ctx *fasthttp.RequestCtx, body any) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	b, err := json.Marshal(body)
	if err != nil {
		RespondError(ctx, fasthttp.StatusInternalServerError, fmt.Errorf("marshal error: %w", err))

		return
	}

	if _, err = ctx.Write(b); err != nil {
		ctx.Logger().Printf("write response error: %v", err)
	}
}

func RespondError(ctx *fasthttp.RequestCtx, code int, err error) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)

	body, e := json.Marshal(map[string]string{"error": err.Error()})
	if e != nil {
		ctx.Logger().Printf("failed to marshal error response: %v", e)

		if _, writeErr := ctx.WriteString(`{"error":"internal server error"}`); writeErr != nil {
			ctx.Logger().Printf("failed to write fallback error response: %v", writeErr)
		}

		return
	}

	if _, writeErr := ctx.Write(body); writeErr != nil {
		ctx.Logger().Printf("write error response: %v", writeErr)
	}
}
