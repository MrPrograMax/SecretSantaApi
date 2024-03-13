package service

import (
	"testApi"
	"testApi/pkg/repository"
)

type ParticipantService struct {
	repo      repository.Participant
	groupRepo repository.Group
}

func NewParticipantService(repo repository.Participant, groupRepo repository.Group) *ParticipantService {
	return &ParticipantService{repo: repo, groupRepo: groupRepo}
}

func (s ParticipantService) Create(groupId int, participant testApi.Participant) (int, error) {
	_, err := s.groupRepo.GetById(groupId)
	if err != nil {
		//group doesn't exist or belongs to user
		return 0, err
	}

	return s.repo.Create(groupId, participant)
}

func (s ParticipantService) Toss(groupId int) ([]testApi.ParticipantDTO, error) {
	return s.repo.Toss(groupId)
}

func (s ParticipantService) GetRecipientInfo(groupId, participantId int) (testApi.RecipientDTO, error) {
	return s.repo.GetRecipientInfo(groupId, participantId)
}

func (s ParticipantService) Delete(groupId, recipientId int) error {
	return s.repo.Delete(groupId, recipientId)
}
