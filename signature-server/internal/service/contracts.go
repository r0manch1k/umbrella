package service

import (
	"crypto/rsa"
	"time"
)

type (
	// SignatureService - сервис для работы с лицензиями.
	SignatureService interface {
		Issue(userID string, duration time.Duration) (string, error)
		Verify(secretPayload string) (string, error)
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
