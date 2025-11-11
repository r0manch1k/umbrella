package signature

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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

	// проверяем лицензию в БД
	lic, err := s.licenseRepo.GetByFingerprint(context.Background(), payload.License, payload.Fingerprint)
	valid := false
	if err == nil && lic != nil {
		if !lic.Activated {
			lic.Activated = true
			_ = s.licenseRepo.Save(context.Background(), *lic)
		}
		valid = true
	}

	resp := struct {
		Valid bool `json:"valid"`
	}{Valid: valid}

	respBytes, _ := json.Marshal(resp)
	hash := sha256.Sum256(respBytes)

	sig, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}
