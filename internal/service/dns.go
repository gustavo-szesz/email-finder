package service

import (
	"net"
	"strings"

	"moremail/email-finder/internal/domain"
)

var _ domain.DNSResult

type DNSService struct{}

func NewDNSService() *DNSService {
	return &DNSService{}
}

func (s *DNSService) Analyze(domainName string) (*domain.DNSResult, error) {
	result := &domain.DNSResult{}

	// MX
	mxRecords, _ := net.LookupMX(domainName)
	for _, mx := range mxRecords {
		result.MX = append(result.MX, mx.Host)
	}

	// TXT (SPF + DMARC)
	txtRecords, _ := net.LookupTXT(domainName)
	for _, txt := range txtRecords {
		if strings.HasPrefix(txt, "v=spf1") {
			result.SPF = txt
		}
	}

	// DMARC usa subdomínio
	dmarcRecords, _ := net.LookupTXT("_dmarc." + domainName)
	for _, txt := range dmarcRecords {
		if strings.HasPrefix(txt, "v=DMARC1") {
			result.DMARC = txt
		}
	}

	return result, nil
}
