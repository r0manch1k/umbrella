package license

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/dto"
)

const randomBytesCapacity = 6

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

// generateAnonymousID — генерируем уникальный идентификатор.
func generateAnonymousID(_ context.Context) string {
	timestamp := time.Now().UTC().Format("20060102150405")

	randomBytes := make([]byte, randomBytesCapacity)
	if _, err := rand.Read(randomBytes); err != nil {
		return timestamp
	}

	return fmt.Sprintf("%s-%s", timestamp, base64.RawURLEncoding.EncodeToString(randomBytes))
}
