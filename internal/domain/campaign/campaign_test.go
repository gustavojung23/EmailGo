package campaign

import "testing"

func TestNewCampaign(t *testing.T) {
	name := "Campaign X"
	content := "Body"
	contacts := []string{"email1@e.com", "email2@e.com"}

	campaign := NewCampaign(name, content, contacts)

	if campaign.ID != "1" {
		t.Errorf("Error ID Email")
	} else if campaign.Name != name {
		t.Errorf("Error Campaign name")
	} else if campaign.Content != content {
		t.Errorf("Error Content")
	} else if len(campaign.Contacts) != len(contacts) {
		t.Errorf("Error contacts")
	}
}
