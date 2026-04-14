package service

import (
	"strings"

	"moremail/email-finder/internal/domain"
)

var providerPatterns = map[domain.EmailProvider][]string{
	domain.ProviderGoogle:    {"google"},
	domain.ProviderMicrosoft: {"outlook", "protection"},
	domain.ProviderZoho:      {"zoho"},
	domain.ProviderProton:    {"proton"},
	domain.ProviderYahoo:     {"yahoo"},
}

func detectProvider(mxRecords []string) domain.EmailProvider {
	for _, mx := range mxRecords {
		mx = strings.ToLower(mx)

		for provider, patterns := range providerPatterns {
			for _, p := range patterns {
				if strings.Contains(mx, p) {
					return provider
				}
			}
		}
	}
	return domain.ProviderUnknown
}
