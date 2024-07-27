package campaign_test

import (
	"emailgo/internal/contract"
	"emailgo/internal/domain/campaign"
	internalerrors "emailgo/internal/internal-errors"
	internalmock "emailgo/internal/test/internalmock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Body Hi!",
		Emails:    []string{"test1@test.com"},
		CreatedBy: "teste@test.com.br",
	}
	campaignPendenting, campaignStarted *campaign.Campaign
	repositoryMock                      *internalmock.CampaignRepositoryMock
	service                             = campaign.ServiceImp{}
)

func setupServiceTest() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPendenting, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func Test_Create_Campaign(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setupServiceTest()
	_, err := service.Create(contract.NewCampaign{})

	assert.False(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPendenting.ID
	})).Return(campaignPendenting, nil)

	campaignReturned, _ := service.GetBy(campaignPendenting.ID)

	assert.Equal(t, campaignPendenting.ID, campaignReturned.ID)
	assert.Equal(t, campaignPendenting.Name, campaignReturned.Name)
	assert.Equal(t, campaignPendenting.Content, campaignReturned.Content)
	assert.Equal(t, campaignPendenting.Status, campaignReturned.Status)
	assert.Equal(t, campaignPendenting.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy("invalid campaign")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problem(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPendenting.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_has_success(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPendenting == campaign
	})).Return(nil)

	err := service.Delete(campaignPendenting.ID)

	assert.Nil(t, err)
}

func Test_Start_ReturnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("campaign invalid")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Start_should_send_email(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)

	emailWasSet := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPendenting.ID {
			emailWasSet = true
		}
		return nil
	}

	service.SendMail = sendMail

	service.Start(campaignPendenting.ID)
	assert.True(t, emailWasSet)
}

func Test_Start_ReturnError_when_func_SendMail_fail(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)

	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignPendenting.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_Return_when_update_to_done(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPendenting.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)
	service.Repository = repositoryMock

	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPendenting.ID)

	assert.Equal(t, campaign.Done, campaignPendenting.Status)
}
