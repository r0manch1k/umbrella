package keypair

import (
	"os"
	"path/filepath"

	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
)

var _ service.KeyPairService = (*Service)(nil)

type Service struct {
	privatePath string
	publicPath  string
	rsaKeyBits  int
}

func New(privatePath, publicPath string, rsaKeyBits int) (*Service, error) {
	err := validatePaths(privatePath, publicPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		privatePath: privatePath,
		publicPath:  publicPath,
		rsaKeyBits:  rsaKeyBits,
	}, nil
}

// validatePaths проверяет корректность путей к ключам.
func validatePaths(privatePath, publicPath string) error {
	if privatePath == "" {
		return exception.ErrEmptyPrivateKeyPath
	}

	if publicPath == "" {
		return exception.ErrEmptyPublicKeyPath
	}

	// Проверяем директорию для приватного ключа
	privateDir := filepath.Dir(privatePath)
	if err := validateDir(privateDir, exception.ErrInvalidPrivateKeyPath); err != nil {
		return err
	}

	// Проверяем директорию для публичного ключа
	publicDir := filepath.Dir(publicPath)
	if err := validateDir(publicDir, exception.ErrInvalidPublicKeyPath); err != nil {
		return err
	}

	// Если файлы уже существуют — убеждаемся, что это не директории
	if fi, err := os.Stat(privatePath); err == nil && fi.IsDir() {
		return exception.ErrInvalidPrivateKeyPath
	}

	if fi, err := os.Stat(publicPath); err == nil && fi.IsDir() {
		return exception.ErrInvalidPublicKeyPath
	}

	return nil
}

// validateDir проверяет, что путь существует и является директорией.
func validateDir(dir string, errToReturn error) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return errToReturn
		}

		return err
	}

	if !info.IsDir() {
		return errToReturn
	}

	return nil
}
