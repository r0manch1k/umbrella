package exception

import "errors"

var (
	ErrInvalidProduct  = errors.New("invalid product")
	ErrLicenseExpired  = errors.New("license expired")
	ErrLicenseNotFound = errors.New("license not found")
)
