package domain

type SMTPResult struct {
	Valid bool   `json:"valid"`
	Host  string `json:"host"`
	Error string `json:"error,omitempty"`
}
