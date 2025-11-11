package exception

import "errors"

var (
	ErrInvalidPrivateKey  = errors.New("invalid private key PEM")
	ErrUnsupportedKeyType = errors.New("unsupported private key type")
)
