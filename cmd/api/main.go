package main

import (
	"moremail/email-finder/internal/api"
	"moremail/email-finder/internal/repository"
	"moremail/email-finder/internal/service"
)

func main() {
	repo := repository.NewMemoryRepository()
	svc := service.NewAnalysisService(repo)
	handler := api.NewHandler(svc)

	r := api.SetupRouter(handler)
	r.Run(":8080")
}
