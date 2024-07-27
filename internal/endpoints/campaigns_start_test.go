package endpoints

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignStart_200(t *testing.T) {
	setupTest()
	campaignId := "xpto"

	service.On("Start", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)

	req, rr := newHttpTest("PATCH", "/", nil)
	req = addParameter(req, "id", campaignId)

	_, status, err := handler.CampaignStart(rr, req)

	assert.Equal(t, 200, status)
	assert.Nil(t, err)
}

func Test_CampaignStart_Err(t *testing.T) {
	setupTest()
	errExpected := errors.New("something wrong")
	service.On("Start", mock.Anything).Return(errExpected)

	req, rr := newHttpTest("PATCH", "/", nil)

	_, _, err := handler.CampaignStart(rr, req)

	assert.Equal(t, errExpected, err)
}
