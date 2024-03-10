package service

import (
	"testApi"
	"testApi/pkg/repository"
)

type Group interface {
	Create(group testApi.Group) (int, error)
	GetAll() ([]testApi.Group, error)
	GetById(groupId int) (testApi.Group, error)
	Delete(groupId int) error
	Update(groupId int, input testApi.UpdateGroupInput) error
}

type Participant interface {
}

type Service struct {
	Group
	Participant
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Group: NewGroupService(repos.Group),
	}
}
