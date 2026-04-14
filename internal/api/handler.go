package api

import (
	"net/http"

	"moremail/email-finder/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.AnalysisService
}

func NewHandler(s *service.AnalysisService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateAnalysis(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.Service.Create(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetAnalysis(c *gin.Context) {
	id := c.Param("id")

	result, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
