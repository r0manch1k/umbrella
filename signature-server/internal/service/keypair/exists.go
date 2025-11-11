package keypair

import (
	"os"
)

func (s *Service) IsExistsPrivateKey() bool {
	_, err := os.Stat(s.privatePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *Service) IsExistsPublicKey() bool {
	_, err := os.Stat(s.publicPath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
