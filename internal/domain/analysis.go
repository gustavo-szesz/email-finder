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

type AnalysisTimings struct {
	TotalMs      int64 `json:"total_ms"`
	DNSMs        int64 `json:"dns_ms"`
	WhoisMs      int64 `json:"whois_ms"`
	SMTPMs       int64 `json:"smtp_ms"`
	CatchAllMs   int64 `json:"catch_all_ms"`
	DisposableMs int64 `json:"disposable_ms"`
	TypoMs       int64 `json:"typo_ms"`
	RelatedMs    int64 `json:"related_ms"`
	RiskMs       int64 `json:"risk_ms"`
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
	Timings  *AnalysisTimings     `json:"timings,omitempty"`
}
