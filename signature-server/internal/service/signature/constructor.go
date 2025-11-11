package signature

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
	"github.com/r0manch1k/umbrella/signature-server/internal/repository"
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
)

const (
	nonceSize = 16
)

var _ service.SignatureService = (*Service)(nil)

type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	product    string

	licenseRepo repository.LicenseRepository
}

func New(privateKeyPath, product string, repo repository.LicenseRepository) (*Service, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		privateKey:  privateKey,
		publicKey:   &privateKey.PublicKey,
		product:     product,
		licenseRepo: repo,
	}, nil
}

// loadPrivateKey считывает приватный ключ из PEM-файла и парсит его.
func loadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	raw, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, exception.ErrInvalidPrivateKey
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}

		return nil, exception.ErrUnsupportedKeyType
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, exception.ErrInvalidPrivateKey
	}

	return rsaKey, nil
}
