package keypair

import (
	"os"
)

const writePrivateKeyPerm os.FileMode = 0o600

func (s *Service) SavePrivateKey(privatePEM []byte) error {
	if err := os.WriteFile(s.privatePath, privatePEM, writePrivateKeyPerm); err != nil {
		return err
	}

	return nil
}

func (s *Service) SavePublicKey(publicPEM []byte) error {
	if err := os.WriteFile(s.publicPath, publicPEM, writePrivateKeyPerm); err != nil {
		return err
	}

	return nil
}
