package keypair

import (
	"os"

	"github.com/r0manch1k/umbrella/signature-server/internal/exception"
)

func (s *Service) GetPublicKey() ([]byte, error) {
	publicPEM, err := os.ReadFile(s.publicPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, exception.ErrPublicKeyNotFound
	}

	return publicPEM, nil
}
