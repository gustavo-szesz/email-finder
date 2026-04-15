package domain

import (
	"time"
)

type DNSResult struct {
	MX       []string `json:"mx"`
	SPF      string   `json:"spf"`
	DMARC    string   `json:"dmarc"`
	Provider string   `json:"provider"`
}

type EmailAnalysis struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Domain    string    `json:"domain"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DNS  *DNSResult  `json:"dns,omitempty"`
	Risk *RiskScore  `json:"risk,omitempty"`
	SMTP *SMTPResult `json:"smtp,omitempty"`

	Disposable *DisposableResult `json:"disposable,omitempty"`
	Typo       *TypoResult       `json:"typo,omitempty"`

	Related  *RelatedEmailsResult `json:"related,omitempty"`
	CatchAll *CatchAllResult      `json:"catch_all,omitempty"`
	Whois    *WhoisResult         `json:"whois,omitempty"`
}
