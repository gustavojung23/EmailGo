package main

import (
	"emailgo/internal/domain/campaign"
	"emailgo/internal/endpoints"
	"emailgo/internal/infrastructure/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db := database.NewDatabase()

	campaignService := campaign.ServiceImp{
		Repository: &database.CampaignRepository{Db: db},
	}

	handler := endpoints.Handler{
		CampaignService: &campaignService,
	}

	r.Post("/campaigns", endpoints.HandlerError(handler.CampaignPost))
	r.Get("/campaigns/{id}", endpoints.HandlerError(handler.CampaignGetById))
	r.Delete("/campaigns/delete/{id}", endpoints.HandlerError(handler.CampaignDelete))

	http.ListenAndServe(":3000", r)
}
