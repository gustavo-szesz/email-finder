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
	requestStart := time.Now()

	a := &domain.EmailAnalysis{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:     email,
		Domain:    util.ExtractDomain(email),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Timings:   &domain.AnalysisTimings{},
	}

	// ======================
	// 1. DNS (primeiro)
	// ======================
	type timedDNSResult struct {
		result *domain.DNSResult
		ms     int64
	}

	dnsChan := make(chan timedDNSResult)

	go func() {
		defer safeRecover(dnsChan)
		started := time.Now()

		result, err := s.DNS.Analyze(a.Domain)
		if err != nil {
			dnsChan <- timedDNSResult{result: nil, ms: time.Since(started).Milliseconds()}
			return
		}
		dnsChan <- timedDNSResult{result: result, ms: time.Since(started).Milliseconds()}
	}()

	dnsTimed := <-dnsChan
	dnsResult := dnsTimed.result
	a.DNS = dnsResult
	a.Timings.DNSMs = dnsTimed.ms

	// ======================
	// 2. Paralelo (WHOIS, SMTP, CatchAll)
	// ======================

	type timedWhoisResult struct {
		result *domain.WhoisResult
		ms     int64
	}
	type timedSMTPResult struct {
		result *domain.SMTPResult
		ms     int64
	}
	type timedCatchAllResult struct {
		result *domain.CatchAllResult
		ms     int64
	}

	whoisChan := make(chan timedWhoisResult)
	smtpChan := make(chan timedSMTPResult)
	catchChan := make(chan timedCatchAllResult)

	// WHOIS
	go func() {
		defer safeRecover(whoisChan)
		started := time.Now()

		result, err := s.Whois.Lookup(a.Domain)
		if err != nil {
			whoisChan <- timedWhoisResult{result: nil, ms: time.Since(started).Milliseconds()}
			return
		}
		whoisChan <- timedWhoisResult{result: result, ms: time.Since(started).Milliseconds()}
	}()

	// SMTP
	go func() {
		defer safeRecover(smtpChan)
		started := time.Now()

		if dnsResult == nil {
			smtpChan <- timedSMTPResult{result: nil, ms: time.Since(started).Milliseconds()}
			return
		}

		smtpChan <- timedSMTPResult{result: s.SMTP.Check(a.Email, a.Domain, dnsResult.MX), ms: time.Since(started).Milliseconds()}
	}()

	// Catch-all
	go func() {
		defer safeRecover(catchChan)
		started := time.Now()

		if dnsResult == nil {
			catchChan <- timedCatchAllResult{result: nil, ms: time.Since(started).Milliseconds()}
			return
		}

		catchChan <- timedCatchAllResult{result: s.CatchAll.Check(a.Domain, dnsResult.MX), ms: time.Since(started).Milliseconds()}
	}()

	// ======================
	// 3. Esperar resultados
	// ======================
	whoisTimed := <-whoisChan
	smtpTimed := <-smtpChan
	catchTimed := <-catchChan

	a.Whois = whoisTimed.result
	a.SMTP = smtpTimed.result
	a.CatchAll = catchTimed.result

	a.Timings.WhoisMs = whoisTimed.ms
	a.Timings.SMTPMs = smtpTimed.ms
	a.Timings.CatchAllMs = catchTimed.ms

	// ======================
	// 4. Sync (rápidos)
	// ======================
	startedDisposable := time.Now()
	a.Disposable = s.Disposable.Check(a.Domain)
	a.Timings.DisposableMs = time.Since(startedDisposable).Milliseconds()

	startedTypo := time.Now()
	a.Typo = s.Typo.Check(a.Domain)
	a.Timings.TypoMs = time.Since(startedTypo).Milliseconds()

	startedRelated := time.Now()
	a.Related = s.Related.Generate(a.Email, a.Domain)
	a.Timings.RelatedMs = time.Since(startedRelated).Milliseconds()

	// ======================
	// 5. Risk
	// ======================
	startedRisk := time.Now()
	a.Risk = s.Risk.Calculate(a.DNS, a.Whois)
	a.Timings.RiskMs = time.Since(startedRisk).Milliseconds()
	a.Timings.TotalMs = time.Since(requestStart).Milliseconds()

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
