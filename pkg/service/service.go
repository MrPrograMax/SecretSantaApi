package service

import (
	"testApi"
	"testApi/pkg/repository"
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

type Service struct {
	Group
	Participant
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Group:       NewGroupService(repos.Group),
		Participant: NewParticipantService(repos.Participant, repos.Group),
	}
}
