package dto

type LicenseVerifyRequest struct {
	SecretPayload string `json:"secret_payload"`
}

type LicenseVerifyResponse struct {
	Signature string `json:"signature"`
}
