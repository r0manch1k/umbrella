package di

import (
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
	"github.com/r0manch1k/umbrella/signature-server/internal/service/keypair"
	"github.com/r0manch1k/umbrella/signature-server/internal/service/signature"
)

type Services struct {
	KeyPair   service.KeyPairService
	Signature service.SignatureService
}

func initServices(cfg *config.Config, repositories Repositories) (Services, error) {
	keypairService, err := keypair.New(
		cfg.Signature.PrivateKeyPath,
		cfg.Signature.PublicKeyPath,
		cfg.Signature.RSAKeyBits,
	)
	if err != nil {
		return Services{}, err
	}

	signatureService, err := signature.New(
		cfg.Signature.PrivateKeyPath,
		cfg.Signature.Product,
		repositories.License,
	)
	if err != nil {
		return Services{}, err
	}

	return Services{
		KeyPair:   keypairService,
		Signature: signatureService,
	}, nil
}
