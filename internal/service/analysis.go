package service

import (
	"fmt"
	"time"

	"moremail/email-finder/internal/domain"
	"moremail/email-finder/internal/repository"
	"moremail/email-finder/internal/util"
)

type AnalysisService struct {
	Repo *repository.MemoryRepository
	DNS  *DNSService
	Risk *RiskService
}

func NewAnalysisService(r *repository.MemoryRepository) *AnalysisService {
	return &AnalysisService{
		Repo: r,
		DNS:  NewDNSService(),
		Risk: NewRiskService(),
	}
}

func (s *AnalysisService) Create(email string) (*domain.EmailAnalysis, error) {
	a := &domain.EmailAnalysis{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:     email,
		Domain:    util.ExtractDomain(email),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 🔥 DNS (uma vez só)
	dnsResult, _ := s.DNS.Analyze(a.Domain)
	a.DNS = dnsResult

	// 🔥 Risk usa o MESMO resultado
	risk := s.Risk.Calculate(dnsResult)
	a.Risk = risk

	a.Status = "done"

	s.Repo.Save(a)
	return a, nil
}

func (s *AnalysisService) GetByID(id string) (*domain.EmailAnalysis, error) {
	return s.Repo.FindByID(id)
}
