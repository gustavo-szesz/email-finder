package domain

type EmailProvider int

const (
	ProviderUnknown EmailProvider = iota
	ProviderGoogle
	ProviderMicrosoft
	ProviderZoho
	ProviderProton
	ProviderYahoo
)

func (p EmailProvider) String() string {
	switch p {
	case ProviderGoogle:
		return "Google"
	case ProviderMicrosoft:
		return "Microsoft"
	case ProviderZoho:
		return "Zoho"
	case ProviderProton:
		return "Proton"
	case ProviderYahoo:
		return "Yahoo"
	default:
		return "Unknown"
	}
}
