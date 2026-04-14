package api

import "github.com/gin-gonic/gin"

func SetupRouter(h *Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/analyze", h.CreateAnalysis)
	router.GET("/analyze/:id", h.GetAnalysis)

	return router
}
