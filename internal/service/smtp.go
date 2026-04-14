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

func (s *SMTPService) Check(email, domainName string, mxRecords []string) *domain.SMTPResult {
	if len(mxRecords) == 0 {
		return &domain.SMTPResult{
			Valid: false,
			Error: "no mx records",
		}
	}

	mxHost := strings.TrimSuffix(mxRecords[0], ".")

	conn, err := net.DialTimeout("tcp", mxHost+":25", 5*time.Second)
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
	}
}
