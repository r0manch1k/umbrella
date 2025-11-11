package exception

import "errors"

var (
	ErrLicenseExpired                       = errors.New("license expired")
	ErrLicenseNotFound                      = errors.New("license not found")
	ErrFailedToVerify                       = errors.New("failed to verify")
	ErrFailedToSaveLicense                  = errors.New("failed to save activated license")
	ErrFailedToSign                         = errors.New("failed to sign")
	ErrLicenseAlreadyActivatedAndNotExpired = errors.New("license already active and not expired")
)
