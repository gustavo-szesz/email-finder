package service

import (
	"strings"

	"moremail/email-finder/internal/domain"
)

type TypoService struct{}

func NewTypoService() *TypoService {
	return &TypoService{}
}

// dominios populares
var commonDomains = []string{
	"gmail.com",
	"outlook.com",
	"yahoo.com",
}

// Comparação
// TODO: Melhorar e ampliar esssa implementação
func isSimilar(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	diff := 0
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff == 1
}

func (s *TypoService) Check(domainName string) *domain.TypoResult {
	domainName = strings.ToLower(domainName)

	for _, d := range commonDomains {
		if isSimilar(domainName, d) {
			return &domain.TypoResult{
				Suspicious: true,
				Target:     d,
			}
		}
	}

	return &domain.TypoResult{
		Suspicious: false,
	}
}
