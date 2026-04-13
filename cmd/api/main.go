package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailAnalysis struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Domain    string    `json:"domain"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"UpdatedAt"`

	//DNS 				*DNSResult			`json:"dns,omitempty"`
	//SMTP 				*SMTPResult 		`json:"smtp,omitempty"`
	//	OSINT     *OSINTResult    `json:"osint,omitempty"`
	///	Breaches  []BreachResult  `json:"breaches,omitempty"`
	//	Risk      *RiskScore      `json:"risk,omitempty"`
}

var analyses []EmailAnalysis

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func analyzeEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	analysis := EmailAnalysis{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Email:     req.Email,
		Domain:    extractDomain(req.Email),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	analyses = append(analyses, analysis)

	c.JSON(201, analysis)
}

func getAnalysis(c *gin.Context) {
	id := c.Param("id")

	for _, a := range analyses {
		if a.ID == id {
			c.JSON(200, a)
			return
		}
	}

	c.JSON(404, gin.H{"error": "not found"})
}

func main() {
	routes := gin.Default()

	routes.POST("/analyze", analyzeEmail)
	routes.GET("/analyze/:id", getAnalysis)

	routes.Run(":8080")
}
