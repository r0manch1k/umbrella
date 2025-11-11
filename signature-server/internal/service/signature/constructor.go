package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
	"github.com/r0manch1k/umbrella/signature-server/internal/repository"
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
)

const (
	nonceSize                      = 16
	privateKeyPEMType  string      = "PRIVATE KEY"
	publicKeyPEMType   string      = "PUBLIC KEY"
	privateKeyFilePerm os.FileMode = 0o600
	publicKeyFilePerm  os.FileMode = 0o644
)

var _ service.SignatureService = (*Service)(nil)

type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	product    string

	licenseRepo repository.LicenseRepository
}

// New создаёт сервис и автоматически генерирует ключи, если их нет.
func New(privateKeyPath, publicKeyPath, product string, privateKeyBits int, repo repository.LicenseRepository) (*Service, error) {
	s := &Service{
		product:     product,
		licenseRepo: repo,
	}

	if !fileExists(privateKeyPath) || !fileExists(publicKeyPath) {
		if err := generateAndSaveKeyPair(privateKeyPath, publicKeyPath, privateKeyBits); err != nil {
			return nil, err
		}
	}

	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	s.privateKey = privateKey
	s.publicKey = &privateKey.PublicKey

	return s, nil
}

// loadPrivateKey считывает приватный ключ из PEM-файла и парсит его.
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, exception.ErrInvalidPrivateKey
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, exception.ErrUnsupportedKeyType
		}

		return rsaKey, nil
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, exception.ErrInvalidPrivateKey
	}

	return rsaKey, nil
}

// generateAndSaveKeyPair генерирует RSA ключи и сохраняет их в файлы.
func generateAndSaveKeyPair(privateKeyPath, publicKeyPath string, privateKeyBits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, privateKeyBits)
	if err != nil {
		return err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  privateKeyPEMType,
		Bytes: privateBytes,
	})

	if err := os.WriteFile(privateKeyPath, privatePEM, privateKeyFilePerm); err != nil {
		return err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  publicKeyPEMType,
		Bytes: pubBytes,
	})

	if err := os.WriteFile(publicKeyPath, pubPEM, publicKeyFilePerm); err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
