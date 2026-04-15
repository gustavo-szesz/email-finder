package service

import (
	"time"

	"moremail/email-finder/internal/domain"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type WhoisService struct{}

func NewWhoisService() *WhoisService {
	return &WhoisService{}
}

func (s *WhoisService) Lookup(domainName string) (*domain.WhoisResult, error) {
	raw, err := whois.Whois(domainName)
	if err != nil {
		return nil, err
	}

	parsed, err := whoisparser.Parse(raw)
	if err != nil {
		return nil, err
	}

	var created time.Time
	var expires time.Time

	if parsed.Domain.CreatedDateInTime != nil {
		created = *parsed.Domain.CreatedDateInTime
	}

	if parsed.Domain.ExpirationDateInTime != nil {
		expires = *parsed.Domain.ExpirationDateInTime
	}

	var age int
	if !created.IsZero() {
		age = int(time.Since(created).Hours() / 24)
	}

	return &domain.WhoisResult{
		Domain:    domainName,
		CreatedAt: created,
		ExpiresAt: expires,
		Registrar: parsed.Registrar.Name,
		AgeInDays: age,
	}, nil
}
