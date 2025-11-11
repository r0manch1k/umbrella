package exception

import "errors"

var (
	ErrEmptyPrivateKeyPath   = errors.New("private key path is empty")
	ErrInvalidPrivateKeyPath = errors.New("invalid private key path")

	ErrEmptyPublicKeyPath   = errors.New("public key path is empty")
	ErrInvalidPublicKeyPath = errors.New("invalid public key path")

	ErrInvalidPrivateKey  = errors.New("invalid private key PEM")
	ErrUnsupportedKeyType = errors.New("unsupported private key type")

	ErrPublicKeyNotFound = errors.New("public key not found")
	ErrReadingPublicKey  = errors.New("error reading public key")
)
