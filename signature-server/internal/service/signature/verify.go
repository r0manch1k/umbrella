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

func (s *Service) Verify(secretPayload string) string {
	data, err := base64.StdEncoding.DecodeString(secretPayload)
	if err != nil {
		return encodeVerifyResponse(false, "")
	}

	var payload struct {
		License     string `json:"license"`
		Fingerprint string `json:"fingerprint"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return encodeVerifyResponse(false, "")
	}

	encryptedLicense, err := base64.StdEncoding.DecodeString(payload.License)
	if err != nil {
		return encodeVerifyResponse(false, "")
	}

	decryptedLicense, err := rsa.DecryptPKCS1v15(rand.Reader, s.privateKey, encryptedLicense)
	if err != nil {
		return encodeVerifyResponse(false, "")
	}

	var licFromClient entity.License
	if err := json.Unmarshal(decryptedLicense, &licFromClient); err != nil {
		return encodeVerifyResponse(false, "")
	}

	licenseDB, err := s.licenseRepo.GetByFingerprint(context.Background(), payload.Fingerprint)
	if err != nil {
		return encodeVerifyResponse(false, "")
	}
	if licenseDB == nil {
		return encodeVerifyResponse(false, "")
	}
	if licenseDB.ExpiresAt.Before(time.Now()) {
		return encodeVerifyResponse(false, "")
	}

	valid := licenseDB.Fingerprint == licFromClient.Fingerprint &&
		licenseDB.Nonce == licFromClient.Nonce &&
		licenseDB.Product == licFromClient.Product

	if valid && !licenseDB.Activated {
		licenseDB.Activated = true
		_ = s.licenseRepo.Save(context.Background(), licenseDB)
	}

	// Формируем подпись для клиента
	licenseJSON, _ := json.Marshal(licenseDB)
	hash := sha256.Sum256(licenseJSON)
	encryptedForClient, _ := rsa.SignPSS(rand.Reader, s.privateKey, crypto.SHA256, hash[:], nil)
	signature := base64.StdEncoding.EncodeToString(encryptedForClient)

	return encodeVerifyResponse(valid, signature)
}

// encodeVerifyResponse - вспомогательная функция для формирования ответа.
func encodeVerifyResponse(valid bool, signature string) string {
	resp := struct {
		Valid     bool   `json:"valid"`
		Signature string `json:"signature,omitempty"`
	}{
		Valid:     valid,
		Signature: signature,
	}

	respBytes, _ := json.Marshal(resp)

	return base64.StdEncoding.EncodeToString(respBytes)
}
