package signature

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
)

func (s *Service) Issue(fingerprint string, duration time.Duration) (string, error) {
	ctx := context.Background()
	now := time.Now()

	existing, err := s.licenseRepo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return "", err
	}

	if existing != nil && existing.Activated && existing.ExpiresAt.After(now) {
		return "", exception.ErrLicenseAlreadyActivatedAndNotExpired
	}

	license := &entity.License{
		Fingerprint: fingerprint,
		Product:     s.product,
		IssuedAt:    now,
		ExpiresAt:   now.Add(duration),
		Nonce:       randomNonce(nonceSize),
		Activated:   false,
	}

	// Сохраняем в базу
	if err := s.licenseRepo.Save(ctx, license); err != nil {
		return "", err
	}

	// сериализация лицензии
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		return "", err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, s.publicKey, licenseJSON)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return encoded, nil
}

func randomNonce(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
