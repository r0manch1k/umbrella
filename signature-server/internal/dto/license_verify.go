package dto

// LicenseVerifyRequest — запрос проверки лицензии.
type LicenseVerifyRequest struct {
	License   string `json:"license"`
	Signature string `json:"signature"`
}

// LicenseVerifyResponse — результат проверки.
type LicenseVerifyResponse struct {
	Signature string `json:"signature"` // зашифрованный ответ
}
