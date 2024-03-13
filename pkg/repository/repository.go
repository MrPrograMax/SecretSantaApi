package repository

import (
	"github.com/jmoiron/sqlx"
	"testApi"
)

type Group interface {
	Create(group testApi.Group) (int, error)
	GetAll() ([]testApi.Group, error)
	GetById(groupId int) (testApi.GroupDTO, error)
	Delete(groupId int) error
	Update(groupId int, input testApi.UpdateGroupInput) error
}

type Participant interface {
	Create(groupId int, item testApi.Participant) (int, error)
	Delete(groupId, participantId int) error
	GetRecipientInfo(groupId, participantId int) (testApi.RecipientDTO, error)
}

type Repository struct {
	Group
	Participant
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Group:       NewGroupPostgres(db),
		Participant: NewParticipantPostgres(db),
	}
}
