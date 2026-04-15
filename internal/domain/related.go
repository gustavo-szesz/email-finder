package domain

type RelatedEmail struct {
	Email      string  `json:"email"`
	Confidence float64 `json:"confidence"`
}

type RelatedEmailsResult struct {
	Emails []RelatedEmail `json:"emails"`
}
