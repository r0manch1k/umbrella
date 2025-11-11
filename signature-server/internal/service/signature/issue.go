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

func (s *Service) Issue(userID, hwFingerprint string, duration time.Duration) (encodedPayload, encodedSignature string, err error) {
	now := time.Now().UTC()
	payload := entity.License{
		UserID:        userID,
		Product:       s.product,
		IssuedAt:      now,
		ExpiresAt:     now.Add(duration),
		HWFingerprint: hwFingerprint,
		Nonce:         randomNonce(nonceSize),
	}

	jb, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	hash := sha256.Sum256(jb)

	sig, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", "", err
	}

	encodedPayload = base64.StdEncoding.EncodeToString(jb)
	encodedSignature = base64.StdEncoding.EncodeToString(sig)

	if s.licenseRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), repositoryTimeout)
		defer cancel()

		if err := s.licenseRepo.Save(ctx, payload); err != nil {
			return "", "", err
		}
	}

	return encodedPayload, encodedSignature, nil
}

func randomNonce(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
