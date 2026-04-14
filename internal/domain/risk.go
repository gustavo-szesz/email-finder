package domain

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type RiskScore struct {
	Score   int       `json:"score"`
	Level   RiskLevel `json:"level"`
	Reasons []string  `json:"reasons"`
}
