package service

import (
	"fmt"
	"time"

	"moremail/email-finder/internal/domain"
	"moremail/email-finder/internal/repository"
	"moremail/email-finder/internal/util"
)

type AnalysisService struct {
	Repo       *repository.MemoryRepository
	DNS        *DNSService
	Risk       *RiskService
	SMTP       *SMTPService
	Disposable *DisposableService
	Typo       *TypoService
}

func NewAnalysisService(r *repository.MemoryRepository) *AnalysisService {
	return &AnalysisService{
		Repo:       r,
		DNS:        NewDNSService(),
		Risk:       NewRiskService(),
		SMTP:       NewSMTPService(),
		Disposable: NewDisposableService(),
		Typo:       NewTypoService(),
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

	// DNS
	dnsChan := make(chan *domain.DNSResult)

	go func() {
		dnsResult, _ := s.DNS.Analyze(a.Domain)
		dnsChan <- dnsResult
	}()

	dnsResult := <-dnsChan
	a.DNS = dnsResult

	// SMTP
	smtpChan := make(chan *domain.SMTPResult)
	go func() {
		smtpChan <- s.SMTP.Check(a.Email, a.Domain, dnsResult.MX)
	}()

	// esperar os dois
	a.SMTP = <-smtpChan

	// Disposable & Typo
	a.Disposable = s.Disposable.Check(a.Domain)
	a.Typo = s.Typo.Check(a.Domain)

	// Risk
	a.Risk = s.Risk.Calculate(dnsResult)

	a.Status = "done"

	s.Repo.Save(a)
	return a, nil
}

func (s *AnalysisService) GetByID(id string) (*domain.EmailAnalysis, error) {
	return s.Repo.FindByID(id)
}
