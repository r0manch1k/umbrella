package license

import (
	"context"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

func (uc *UseCase) Issue(ctx context.Context, req dto.LicenseIssueRequest) (dto.LicenseIssueResponse, error) {
	duration := time.Duration(req.DurationHours) * time.Hour

	// Генерируем userID, если пустой
	userID := req.UserID
	if userID == "" {
		userID = generateAnonymousID(ctx)
	}

	license, signature, err := uc.signer.Issue(userID, req.HWFingerprint, duration)
	if err != nil {
		return dto.LicenseIssueResponse{}, err
	}

	return dto.LicenseIssueResponse{
		License:   license,
		Signature: signature,
	}, nil
}

// generateAnonymousID — генерируем идентификатор по ip/device/time.
func generateAnonymousID(_ context.Context) string {
	// TODO: можно использовать хэш ip+user-agent+timestamp
	return time.Now().UTC().Format("20060102150405")
}
