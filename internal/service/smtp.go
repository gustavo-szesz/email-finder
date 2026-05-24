package service

import (
	"net"
	"net/smtp"
	"strings"
	"time"

	"moremail/email-finder/internal/domain"
)

type SMTPService struct{}

func NewSMTPService() *SMTPService {
	return &SMTPService{}
}

func (s *SMTPService) checkSingle(email, mxHost string) *domain.SMTPResult {
	conn, err := net.DialTimeout("tcp", mxHost+":25", 3*time.Second)
	if err != nil {
		return &domain.SMTPResult{
			Valid: false,
			Host:  mxHost,
			Error: err.Error(),
		}
	}

	client, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return &domain.SMTPResult{
			Valid: false,
			Host:  mxHost,
			Error: err.Error(),
		}
	}

	defer client.Close()

	client.Hello("localhost")
	client.Mail("test@example.com")

	err = client.Rcpt(email)
	if err != nil {
		return &domain.SMTPResult{
			Valid: false,
			Host:  mxHost,
			Error: err.Error(),
		}
	}

	return &domain.SMTPResult{
		Valid: true,
		Host:  mxHost,
		Error: "provider blocked smtp validation",
	}
}

func (s *SMTPService) Check(email, domainName string, mxRecords []string) *domain.SMTPResult {
	if len(mxRecords) == 0 {
		return &domain.SMTPResult{
			Valid: false,
			Error: "no mx records",
		}
	}

	resultChan := make(chan *domain.SMTPResult, len(mxRecords))

	for _, mx := range mxRecords {
		mxHost := strings.TrimSuffix(mx, ".")

		go func(host string) {
			res := s.checkSingle(email, host)
			resultChan <- res
		}(mxHost)
	}

	// timeout global
	timeout := time.After(5 * time.Second)

	for i := 0; i < len(mxRecords); i++ {
		select {
		case res := <-resultChan:
			if res.Valid {
				return res // 🔥 primeiro válido ganha
			}
		case <-timeout:
			return &domain.SMTPResult{
				Valid: false,
				Error: "timeout",
			}
		}
	}

	return &domain.SMTPResult{
		Valid: false,
		Error: "all mx failed",
	}
}
