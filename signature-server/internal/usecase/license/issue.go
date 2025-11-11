package license

import (
	"context"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

func (uc *UseCase) Issue(ctx context.Context, req dto.LicenseIssueRequest) (dto.LicenseIssueResponse, error) {
	if req.Fingerprint == "" {
		req.Fingerprint = generateAnonymousID(ctx)
	}

	license, err := uc.signer.Issue(req.Fingerprint, time.Duration(req.DurationHours)*time.Hour)
	if err != nil {
		return dto.LicenseIssueResponse{}, err
	}

	return dto.LicenseIssueResponse{License: license}, nil
}

// generateAnonymousID — генерируем идентификатор по ip/device/time.
func generateAnonymousID(_ context.Context) string {
	return time.Now().UTC().Format("20060102150405")
}
