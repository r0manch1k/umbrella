package di

import (
	"github.com/r0manch1k/umbrella/signature-server/config"
	"github.com/r0manch1k/umbrella/signature-server/internal/service"
	"github.com/r0manch1k/umbrella/signature-server/internal/service/signature"
)

type Services struct {
	Signature service.SignatureService
}

func initServices(cfg *config.Config, repositories Repositories) (Services, error) {
	signatureService, err := signature.New(
		cfg.Signature.PrivateKeyPath,
		cfg.Signature.PublicKeyPath,
		cfg.Signature.Product,
		cfg.Signature.RSAKeyBits,
		repositories.License,
	)
	if err != nil {
		return Services{}, err
	}

	return Services{
		Signature: signatureService,
	}, nil
}
