package signature

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
)

func (s *Service) Issue(fingerprint string, duration time.Duration) (string, error) {
	now := time.Now().UTC()
	license := &entity.License{
		Fingerprint: fingerprint,
		Product:     s.product,
		IssuedAt:    now,
		ExpiresAt:   now.Add(duration),
		Nonce:       randomNonce(nonceSize),
		Activated:   false,
	}

	// Сохраняем лицензию и проверяем ошибку
	if s.licenseRepo != nil {
		if err := s.licenseRepo.Save(context.Background(), license); err != nil {
			return "", err
		}
	}

	// Мерджим JSON и проверяем ошибку
	jb, err := json.Marshal(license)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jb)

	sig, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(sig)

	return encoded, nil
}

func randomNonce(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
