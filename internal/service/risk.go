package service

import "moremail/email-finder/internal/domain"

type RiskService struct{}

func NewRiskService() *RiskService {
	return &RiskService{}
}

func (s *RiskService) Calculate(dns *domain.DNSResult) *domain.RiskScore {
	score := 0
	var reasons []string

	// sem MX
	if len(dns.MX) == 0 {
		score += 50
		reasons = append(reasons, "No MX records")
	}

	// sem SPF
	if dns.SPF == "" {
		score += 20
		reasons = append(reasons, "No SPF record")
	}

	if dns.DMARC == "" {
		score += 15
		reasons = append(reasons, "No DMARC record")
	}
	// provider desconhecido
	if dns.Provider == "Unknown" {
		score += 15
		reasons = append(reasons, "Unknown provider")
	}

	// determine level
	level := domain.RiskLow
	if score >= 60 {
		level = domain.RiskHigh
	} else if score >= 30 {
		level = domain.RiskMedium
	}

	return &domain.RiskScore{
		Score:   score,
		Level:   level,
		Reasons: reasons,
	}
}
