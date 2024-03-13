package service

import (
	"testApi"
	"testApi/pkg/repository"
)

type GroupService struct {
	repo repository.Group
}

func NewGroupService(repo repository.Group) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) Create(group testApi.Group) (int, error) {
	return s.repo.Create(group)
}

func (s *GroupService) GetAll() ([]testApi.Group, error) {
	return s.repo.GetAll()
}

func (s *GroupService) GetById(groupId int) (testApi.GroupDTO, error) {
	return s.repo.GetById(groupId)
}

func (s *GroupService) Delete(groupId int) error {
	return s.repo.Delete(groupId)
}

func (s *GroupService) Update(groupId int, input testApi.UpdateGroupInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(groupId, input)
}
