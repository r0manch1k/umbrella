package dto

// LicenseIssueRequest — запрос для выдачи лицензии.
type LicenseIssueRequest struct {
	UserID        string `json:"user_id,omitempty"` // опционально
	DurationHours int    `json:"duration_hours"`    // срок действия лицензии
	HWFingerprint string `json:"hw_fingerprint,omitempty"`
}

// LicenseIssueResponse — ответ сервера при выдаче лицензии.
type LicenseIssueResponse struct {
	License   string `json:"license"`   // base64 payload
	Signature string `json:"signature"` // base64 подпись
}
