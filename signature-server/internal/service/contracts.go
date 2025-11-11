package service

import (
	"crypto/rsa"
	"time"
)

type (
	// SignatureService - сервис для работы с лицензиями.
	SignatureService interface {
		Issue(userID, hwFingerprint string, duration time.Duration) (payload, signature string, err error)
		Verify(encPayload, encSig string) (string, error)
		Sign(payload []byte) (string, error)
	}

	// KeyPairService - сервис, который.
	KeyPairService interface {
		GeneratePrivateKey() (*rsa.PrivateKey, []byte, error)
		GeneratePublicKey(privateKey *rsa.PrivateKey) ([]byte, error)
		SavePrivateKey(privatePEM []byte) error
		SavePublicKey(publicPEM []byte) error
		IsExistsPrivateKey() bool
		IsExistsPublicKey() bool
		GetPublicKey() ([]byte, error)
	}
)
