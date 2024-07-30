package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	content   = "Body Hi!"
	contacts  = []string{"email1@e.com", "email2@e.com"}
	createdBy = "teste@teste.com.br"

	fake = faker.New()

	campaignNewCampaign *Campaign
)

func setupNewCampaign() {
	campaignNewCampaign, _ = NewCampaign(name, content, contacts, createdBy)
}

func Test_NewCampaign_CreateCampaign(t *testing.T) {
	setupNewCampaign()

	assert.Equal(t, campaignNewCampaign.Name, name)
	assert.Equal(t, campaignNewCampaign.Content, content)
	assert.Equal(t, len(campaignNewCampaign.Contacts), len(contacts))
	assert.Equal(t, campaignNewCampaign.CreatedBy, createdBy)
}

func Test_NewCapaign_IDIsNotNil(t *testing.T) {
	setupNewCampaign()

	assert.NotNil(t, campaignNewCampaign.ID)
}

func Test_NewCampaign_MustStatusStartWithPending(t *testing.T) {
	setupNewCampaign()

	assert.Equal(t, Pending, campaignNewCampaign.Status)
}

func Test_NewCapaign_CreatedOnMustBeNow(t *testing.T) {
	now := time.Now().Add(-time.Minute)

	setupNewCampaign()

	assert.Greater(t, campaignNewCampaign.CreatedOn, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {

	_, err := NewCampaign("", content, contacts, createdBy)

	assert.Equal(t, "name is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, createdBy)

	assert.Equal(t, "name is required with max 24", err.Error())
}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	_, err := NewCampaign(name, "", contacts, createdBy)

	assert.Equal(t, "content is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts, createdBy)

	assert.Equal(t, "content is required with max 1024", err.Error())
}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {
	_, err := NewCampaign(name, content, nil, createdBy)

	assert.Equal(t, "contacts is required with min 1", err.Error())
}

func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	_, err := NewCampaign(name, content, []string{"email_invalid"}, createdBy)

	assert.Equal(t, "email is invalid", err.Error())
}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {
	_, err := NewCampaign(name, content, contacts, "")

	assert.Equal(t, "createdby is invalid", err.Error())
}

func Test_Done_ChangeStatus(t *testing.T) {
	setupNewCampaign()

	campaignNewCampaign.Done()

	assert.Equal(t, Done, campaignNewCampaign.Status)
}

func Test_Start_ChangeStatus(t *testing.T) {
	setupNewCampaign()

	campaignNewCampaign.Started()

	assert.Equal(t, Started, campaignNewCampaign.Status)
}

func Test_Cancel_ChangeStatus(t *testing.T) {
	setupNewCampaign()

	campaignNewCampaign.Cancel()

	assert.Equal(t, Canceled, campaignNewCampaign.Status)
}

func Test_Delete_ChangeStatus(t *testing.T) {
	setupNewCampaign()

	campaignNewCampaign.Delete()

	assert.Equal(t, Deleted, campaignNewCampaign.Status)
}

func Test_Fail_ChangeStatus(t *testing.T) {
	setupNewCampaign()

	campaignNewCampaign.Fail()

	assert.Equal(t, Fail, campaignNewCampaign.Status)
}
