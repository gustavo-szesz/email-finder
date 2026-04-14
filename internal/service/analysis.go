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
	SMTP *SMTPService
}

func NewAnalysisService(r *repository.MemoryRepository) *AnalysisService {
	return &AnalysisService{
		Repo: r,
		DNS:  NewDNSService(),
		Risk: NewRiskService(),
		SMTP: NewSMTPService(),
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

	//  DNS
	dnsResult, _ := s.DNS.Analyze(a.Domain)
	a.DNS = dnsResult
	// SMTP
	smtpResult := s.SMTP.Check(a.Email, a.Domain, dnsResult.MX)
	a.SMTP = smtpResult

	//  Risk
	risk := s.Risk.Calculate(dnsResult)
	a.Risk = risk

	a.Status = "done"

	s.Repo.Save(a)
	return a, nil
}

func (s *AnalysisService) GetByID(id string) (*domain.EmailAnalysis, error) {
	return s.Repo.FindByID(id)
}
