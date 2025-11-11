package entity

import "time"

type License struct {
	Fingerprint string    `json:"fingerprint"`
	Product     string    `json:"product"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Nonce       string    `json:"nonce"`
	Activated   bool      `json:"activated"`
}
