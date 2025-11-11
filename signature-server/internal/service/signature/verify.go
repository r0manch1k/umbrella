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

	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
)

func (s *Service) Verify(secretPayload string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(secretPayload)
	if err != nil {
		return "", err
	}

	var payload struct {
		License     string `json:"license"`
		Fingerprint string `json:"fingerprint"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}

	license, err := s.licenseRepo.GetByFingerprint(context.Background(), payload.Fingerprint)
	if err != nil {
		return "", exception.ErrFailedToVerify
	}

	if license == nil {
		return "", exception.ErrLicenseNotFound
	}

	if license.ExpiresAt.Before(time.Now().UTC()) {
		return "", exception.ErrLicenseExpired
	}

	// Если лицензия ещё не активирована — активируем
	if !license.Activated {
		license.Activated = true
		if err := s.licenseRepo.Save(context.Background(), license); err != nil {
			return "", exception.ErrFailedToSaveLicense
		}
	}

	// формируем ответ
	response := struct {
		Valid bool `json:"valid"`
	}{Valid: true}

	respBytes, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(respBytes)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", exception.ErrFailedToSign
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
