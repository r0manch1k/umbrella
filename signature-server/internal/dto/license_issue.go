package dto

type LicenseIssueRequest struct {
	Fingerprint   string `json:"fingerprint,omitempty"`
	DurationHours int    `json:"duration_hours"`
}

type LicenseIssueResponse struct {
	License string `json:"license"`
}
