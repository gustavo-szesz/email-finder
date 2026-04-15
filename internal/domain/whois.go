package domain

import "time"

type WhoisResult struct {
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Registrar string    `json:"registrar"`
	AgeInDays int       `json:"age_in_days"`
}
