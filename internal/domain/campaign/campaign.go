package campaign

import (
	internalerrors "emailgo/internal/internal-errors"
	"time"

	"github.com/rs/xid"
)

const (
	Pending  = "Pending"
	Started  = "Started"
	Done     = "Done"
	Canceled = "Canceled"
	Deleted  = "Deleted"
	Fail     = "Fail"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignId string `gorm:"size:50"`
}

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    string    `gorm:"size:20"`
	CreatedBy string    `validate:"email" gorm:"size:50"`
}

func (c *Campaign) Done() {
	c.Status = Done
}

func (c *Campaign) Cancel() {
	c.Status = Canceled
}

func (c *Campaign) Delete() {
	c.Status = Deleted
}

func (c *Campaign) Fail() {
	c.Status = Fail
}

func (c *Campaign) Started() {
	c.Status = Started
}

func NewCampaign(name string, content string, emails []string, createdBy string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))
	for index, email := range emails {
		contacts[index].Email = email
		contacts[index].ID = xid.New().String()
	}

	campaing := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		CreatedOn: time.Now(),
		Content:   content,
		Contacts:  contacts,
		Status:    Pending,
		CreatedBy: createdBy,
	}
	err := internalerrors.ValidateStruct(campaing)
	if err == nil {
		return campaing, nil
	} else {
		return nil, err
	}
}
