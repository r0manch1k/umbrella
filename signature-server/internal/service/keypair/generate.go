package keypair

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	privateKeyPEMType string = "PRIVATE KEY"
	publicKeyPEMType  string = "PUBLIC KEY"
)

func (s *Service) GeneratePrivateKey() (*rsa.PrivateKey, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, s.rsaKeyBits)
	if err != nil {
		return nil, nil, err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  privateKeyPEMType,
		Bytes: privateBytes,
	})

	return privateKey, privatePEM, nil
}

func (s *Service) GeneratePublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  publicKeyPEMType,
		Bytes: publicBytes,
	})

	return publicPEM, nil
}
