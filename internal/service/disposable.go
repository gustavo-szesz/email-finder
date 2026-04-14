package service

import (
	"strings"

	"moremail/email-finder/internal/domain"
)

type DisposableService struct{}

func NewDisposableService() *DisposableService {
	return &DisposableService{}
}

// lista inicial (pode crescer depois)
var disposableDomains = map[string]string{
	"10minutemail.com":  "10minutemail",
	"mailinator.com":    "mailinator",
	"guerrillamail.com": "guerrillamail",
	"tempmail.com":      "tempmail",
}

func (s *DisposableService) Check(domainName string) *domain.DisposableResult {
	domainName = strings.ToLower(domainName)

	if provider, ok := disposableDomains[domainName]; ok {
		return &domain.DisposableResult{
			IsDisposable: true,
			Provider:     provider,
		}
	}

	return &domain.DisposableResult{
		IsDisposable: false,
	}
}
