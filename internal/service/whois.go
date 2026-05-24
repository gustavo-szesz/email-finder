package service

import (
	"fmt"
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
	type lookupResult struct {
		result *domain.WhoisResult
		err    error
	}

	ch := make(chan lookupResult, 1)

	go func() {
		raw, err := whois.Whois(domainName)
		if err != nil {
			ch <- lookupResult{result: nil, err: err}
			return
		}

		parsed, err := whoisparser.Parse(raw)
		if err != nil {
			ch <- lookupResult{result: nil, err: err}
			return
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

		ch <- lookupResult{
			result: &domain.WhoisResult{
				Domain:    domainName,
				CreatedAt: created,
				ExpiresAt: expires,
				Registrar: parsed.Registrar.Name,
				AgeInDays: age,
			},
			err: nil,
		}
	}()

	select {
	case out := <-ch:
		return out.result, out.err
	case <-time.After(1200 * time.Millisecond):
		return nil, fmt.Errorf("whois timeout")
	}
}
