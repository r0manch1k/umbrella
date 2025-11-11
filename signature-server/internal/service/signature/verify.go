package signature

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/r0manch1k/umbrella/signature-server/internal/entity"
	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
)

func (s *Service) Verify(encPayload, encSig string) (string, error) {
	// Декодируем и проверяем исходную лицензию
	jb, err := base64.StdEncoding.DecodeString(encPayload)
	if err != nil {
		return "", err
	}

	sig, err := base64.StdEncoding.DecodeString(encSig)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(jb)
	if err := rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hash[:], sig); err != nil {
		return "", exception.ErrInvalidSignature
	}

	var payload entity.License
	if err := json.Unmarshal(jb, &payload); err != nil {
		return "", err
	}

	valid := !(payload.Product != s.product || time.Now().UTC().After(payload.ExpiresAt))

	// Проверяем наличие лицензии в БД
	if s.licenseRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		stored, err := s.licenseRepo.GetByUserAndFingerprint(ctx, payload.UserID, payload.HWFingerprint)
		if err != nil || stored == nil {
			valid = false
		}
	}

	// Формируем payload ответа
	responsePayload := struct {
		Valid bool `json:"valid"`
	}{
		Valid: valid,
	}
	plain, _ := json.Marshal(responsePayload)

	// Сначала шифруем секретным AES методом
	secretEncrypted, err := encryptSecretAES(plain)
	if err != nil {
		return "", err
	}

	// Затем подписываем RSA
	finalSig, err := s.Sign([]byte(secretEncrypted))
	if err != nil {
		return "", err
	}

	return finalSig, nil
}
