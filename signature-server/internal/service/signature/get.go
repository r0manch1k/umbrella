package signature

import (
	"crypto/x509"
	"encoding/pem"
)

func (s *Service) GetPublicKey() ([]byte, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(s.publicKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  publicKeyPEMType,
		Bytes: pubBytes,
	}

	return pem.EncodeToMemory(block), nil
}
