package license

import (
	"context"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

func (uc *UseCase) Verify(_ context.Context, req dto.LicenseVerifyRequest) (dto.LicenseVerifyResponse, error) {
	sig, err := uc.signer.Verify(req.License, req.Signature)
	if err != nil {
		return dto.LicenseVerifyResponse{}, err
	}

	return dto.LicenseVerifyResponse{
		Signature: sig,
	}, nil
}
