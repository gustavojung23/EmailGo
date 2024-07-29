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

func setupSendEmailTest(err error) {
	sendMail := func(campaign *campaign.Campaign) error {
		return err
	}
	service.SendMail = sendMail
}

func Test_Create_RequestIsValid_IdIsNotNil(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_RequestIsNotValid_ErrInternal(t *testing.T) {
	setupServiceTest()
	_, err := service.Create(contract.NewCampaign{})

	assert.False(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_RequestIsValid_CallRepository(t *testing.T) {
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

func Test_Create_ErrorOnRepository_ErrInternal(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_CampaignExists_CampaignSaved(t *testing.T) {
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

func Test_GetById_ErrorOnRepository_ErrInternal(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy("invalid campaign")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_CampaignWasNotFound_ErrRecordNotFound(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_CampaignIsNotPending_Err(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ErrorOnRepository_ErrInternal(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPendenting.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_CampaignWasDeleted_Nil(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPendenting == campaign
	})).Return(nil)

	err := service.Delete(campaignPendenting.ID)

	assert.Nil(t, err)
}

func Test_Start_CamapaignWasNotFound_ErrRecordNotFound(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("campaign invalid")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_CampaignIsNotPending_Err(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Start_CampaignWasUpdated_StatusIsStarted(t *testing.T) {
	setupServiceTest()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPendenting, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPendenting.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Started
	})).Return(nil)

	setupSendEmailTest(nil)

	service.Start(campaignPendenting.ID)

	assert.Equal(t, campaign.Started, campaignPendenting.Status)
}

func Test_SendEmailUpdateStatus_WhenFail_StatusIsFail(t *testing.T) {
	setupServiceTest()
	setupSendEmailTest(errors.New("error to send email"))
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPendenting.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Fail
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPendenting)

	repositoryMock.AssertExpectations(t)
}

func Test_SendEmailUpdateStatus_WhenSuccess_StatusIsDone(t *testing.T) {
	setupServiceTest()
	setupSendEmailTest(nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPendenting.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPendenting)

	repositoryMock.AssertExpectations(t)
}
