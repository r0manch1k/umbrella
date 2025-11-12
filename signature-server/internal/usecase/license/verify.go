package license

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

func (uc *UseCase) Verify(_ context.Context, req dto.LicenseVerifyRequest) (dto.LicenseVerifyResponse, error) {
	sig := uc.signer.Verify(req.SecretPayload)

	return dto.LicenseVerifyResponse{Signature: sig}, nil
}
