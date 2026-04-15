package service

import (
	"fmt"
	"math/rand"
	"time"

	"moremail/email-finder/internal/domain"
)

type CatchAllService struct {
	SMTP *SMTPService
}

func NewCatchAllService(smtp *SMTPService) *CatchAllService {
	return &CatchAllService{SMTP: smtp}
}

// gerar email aleatorio
func randomEmail(domainName string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("noexist_%d@%s", rand.Intn(99999999), domainName)
}

func (s *CatchAllService) Check(domainName string, mxRecords []string) *domain.CatchAllResult {
	testEmail := randomEmail((domainName))

	result := s.SMTP.Check(testEmail, domainName, mxRecords)

	isCatchAll := result.Valid

	return &domain.CatchAllResult{
		IsCatchAll: isCatchAll,
		TestEmail:  testEmail,
	}
}
