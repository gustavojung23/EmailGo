package endpoints

import (
	"emailgo/internal/contract"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	body = contract.NewCampaignRequest{
		Name:    "teste",
		Content: "Hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
	createdByExpected = "teste@teste.com.br"
)

func Test_CampaignsPost_201(t *testing.T) {
	setupTest()

	service.On("Create", mock.MatchedBy(func(request contract.NewCampaignRequest) bool {
		if request.Name == body.Name && request.Content == body.Content && request.CreatedBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("312312", nil)

	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)
	_, status, err := handler.CampaignPost(rr, req)

	assert.Equal(t, 201, status)
	assert.Nil(t, err)
}

func Test_CampaignsPost_Err(t *testing.T) {
	setupTest()

	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))

	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, _, err := handler.CampaignPost(rr, req)

	assert.NotNil(t, err)
}
