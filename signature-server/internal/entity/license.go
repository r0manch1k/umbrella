package entity

import "time"

type License struct {
	UserID        string    `json:"user_id"`
	Product       string    `json:"product"`
	HWFingerprint string    `json:"hw_fingerprint"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	Nonce         string    `json:"nonce"`
}
