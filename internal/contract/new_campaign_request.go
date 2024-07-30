package contract

type NewCampaignRequest struct {
	Name      string
	Content   string
	Emails    []string
	CreatedBy string
}
