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
	Related    *RelatedService
	CatchAll   *CatchAllService
	Whois      *WhoisService
}

func NewAnalysisService(r *repository.MemoryRepository) *AnalysisService {
	return &AnalysisService{
		Repo:       r,
		DNS:        NewDNSService(),
		Risk:       NewRiskService(),
		SMTP:       NewSMTPService(),
		Disposable: NewDisposableService(),
		Typo:       NewTypoService(),
		Related:    NewRelatedService(),
		CatchAll:   NewCatchAllService(NewSMTPService()),
		Whois:      NewWhoisService(),
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

	// ======================
	// 1. DNS (primeiro)
	// ======================
	dnsChan := make(chan *domain.DNSResult)

	go func() {
		defer safeRecover(dnsChan)

		result, err := s.DNS.Analyze(a.Domain)
		if err != nil {
			dnsChan <- nil
			return
		}
		dnsChan <- result
	}()

	dnsResult := <-dnsChan
	a.DNS = dnsResult

	// ======================
	// 2. Paralelo (WHOIS, SMTP, CatchAll)
	// ======================

	whoisChan := make(chan *domain.WhoisResult)
	smtpChan := make(chan *domain.SMTPResult)
	catchChan := make(chan *domain.CatchAllResult)

	// WHOIS
	go func() {
		defer safeRecover(whoisChan)

		result, err := s.Whois.Lookup(a.Domain)
		if err != nil {
			whoisChan <- nil
			return
		}
		whoisChan <- result
	}()

	// SMTP
	go func() {
		defer safeRecover(smtpChan)

		if dnsResult == nil {
			smtpChan <- nil
			return
		}

		smtpChan <- s.SMTP.Check(a.Email, a.Domain, dnsResult.MX)
	}()

	// Catch-all
	go func() {
		defer safeRecover(catchChan)

		if dnsResult == nil {
			catchChan <- nil
			return
		}

		catchChan <- s.CatchAll.Check(a.Domain, dnsResult.MX)
	}()

	// ======================
	// 3. Esperar resultados
	// ======================
	a.Whois = <-whoisChan
	a.SMTP = <-smtpChan
	a.CatchAll = <-catchChan

	// ======================
	// 4. Sync (rápidos)
	// ======================
	a.Disposable = s.Disposable.Check(a.Domain)
	a.Typo = s.Typo.Check(a.Domain)
	a.Related = s.Related.Generate(a.Email, a.Domain)

	// ======================
	// 5. Risk
	// ======================
	a.Risk = s.Risk.Calculate(a.DNS, a.Whois)

	a.Status = "done"

	s.Repo.Save(a)
	return a, nil
}

func (s *AnalysisService) GetByID(id string) (*domain.EmailAnalysis, error) {
	return s.Repo.FindByID(id)
}

func safeRecover[T any](ch chan T) {
	if r := recover(); r != nil {
		var zero T
		ch <- zero
	}
}
