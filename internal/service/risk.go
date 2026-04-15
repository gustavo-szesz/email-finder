package service

import (
	"moremail/email-finder/internal/domain"
)

type RiskService struct{}

func NewRiskService() *RiskService {
	return &RiskService{}
}

func (s *RiskService) Calculate(
	dns *domain.DNSResult,
	whois *domain.WhoisResult,
) *domain.RiskScore {
	score := 0
	var reasons []string

	// DNS
	if dns.SPF == "" {
		score += 20
		reasons = append(reasons, "missing SPF")
	}

	if dns.DMARC == "" {
		score += 20
		reasons = append(reasons, "missing DMARC")
	}

	// WHOIS
	if whois != nil && !whois.CreatedAt.IsZero() {
		if whois.AgeInDays < 30 {
			score += 50
			reasons = append(reasons, "domain very new")
		} else if whois.AgeInDays < 180 {
			score += 30
			reasons = append(reasons, "domain relatively new")
		}
	}

	// LEVEL
	level := domain.RiskLevel("low")

	if score > 70 {
		level = domain.RiskLevel("high")
	} else if score > 40 {
		level = domain.RiskLevel("medium")
	}

	return &domain.RiskScore{
		Score:   score,
		Level:   level,
		Reasons: reasons,
	}
}
