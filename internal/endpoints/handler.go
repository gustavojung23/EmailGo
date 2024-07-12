package endpoints

import "emailgo/internal/domain/campaign"

type Handler struct {
	CampaignService campaign.Service
}
