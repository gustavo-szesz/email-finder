package service

import (
	"fmt"
	"strings"

	"moremail/email-finder/internal/domain"
)

type RelatedService struct{}

func NewRelatedService() *RelatedService {
	return &RelatedService{}
}

// padrões comuns de emails corporativos
var commonUsers = []string{
	"admin",
	"contato",
	"contact",
	"support",
	"suporte",
	"sales",
	"finance",
	"financeiro",
	"info",
}

// extrai o nome antes do @
func extractName(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func (s *RelatedService) Generate(email, domainName string) *domain.RelatedEmailsResult {
	var results []domain.RelatedEmail

	name := extractName(email)

	// padrões comuns
	for _, user := range commonUsers {
		results = append(results, domain.RelatedEmail{
			Email:      fmt.Sprintf("%s@%s", user, domainName),
			Confidence: estimateConfidence(user),
		})
	}

	// variações do usuário
	if name != "" {
		variations := []string{
			name,
			name + "123",
		}

		for _, v := range variations {
			results = append(results, domain.RelatedEmail{
				Email:      fmt.Sprintf("%s@%s", v, domainName),
				Confidence: estimateConfidence(v),
			})
		}
	}

	return &domain.RelatedEmailsResult{
		Emails: results,
	}
}

func estimateConfidence(user string) float64 {
	switch user {
	case "admin", "contato", "contact", "support", "suporte":
		return 0.9

	case "info", "sales", "finance", "financeiro":
		return 0.8

	default:
		if len(user) <= 5 {
			return 0.8
		}
		return 0.4
	}
}
