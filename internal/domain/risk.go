package domain

type RiskLevel string

const (
	Low    RiskLevel = "low"
	Medium RiskLevel = "medium"
	High   RiskLevel = "high"
)

type RiskScore struct {
	Score   int       `json:"score"`
	Level   RiskLevel `json:"level"`
	Reasons []string  `json:"reasons"`
}
