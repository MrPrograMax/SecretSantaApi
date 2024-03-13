package testApi

import "errors"

type Group struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Description string `json:"description" db:"description"`
}

type GroupDTO struct {
	Group
	Participants []ParticipantDTO
}

type GroupsList struct {
	Id            int
	GroupId       int
	ParticipantId int
}

type Participant struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Wish        string `json:"wish" db:"wish"`
	RecipientId int    `json:"recipientId" db:"recipient_id"`
}

type ParticipantDTO struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Wish      string        `json:"wish"`
	Recipient *RecipientDTO `json:"recipient"`
}

type RecipientDTO struct {
	Id   int    `json:"recipient_id"`
	Name string `json:"recipient_name"`
	Wish string `json:"recipient_wish"`
}

type UpdateGroupInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (i UpdateGroupInput) Validate() error {
	if i.Name == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
