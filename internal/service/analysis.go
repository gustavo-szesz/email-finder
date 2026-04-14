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
}

func NewAnalysisService(r *repository.MemoryRepository) *AnalysisService {
	return &AnalysisService{
		Repo: r,
		DNS:  NewDNSService(),
	}
}

func (s *AnalysisService) Create(email string) (*domain.EmailAnalysis, error) {
	a := &domain.EmailAnalysis{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:     email,
		Domain:    util.ExtractDomain(email),
		DNS:       nil,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dnsResult, _ := s.DNS.Analyze(a.Domain)
	a.DNS = dnsResult

	a.Status = "done"

	s.Repo.Save(a)
	return a, nil
}

func (s *AnalysisService) GetByID(id string) (*domain.EmailAnalysis, error) {
	return s.Repo.FindByID(id)
}
