package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"testApi"
)

type ParticipantPostgres struct {
	db *sqlx.DB
}

func NewParticipantPostgres(db *sqlx.DB) *ParticipantPostgres {
	return &ParticipantPostgres{db: db}
}

func (r *ParticipantPostgres) Create(groupId int, item testApi.Participant) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var participantId int
	createParticipantQuery := fmt.Sprintf("INSERT INTO %s (name, wish, recipient_id) VALUES ($1, $2, null) RETURNING id", participantsTable)

	row := tx.QueryRow(createParticipantQuery, item.Name, item.Wish)
	err = row.Scan(&participantId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createGroupParticipantsQuery := fmt.Sprintf("INSERT INTO %s (group_id, participant_id) VALUES ($1, $2)", groupsParticipantsListsTable)
	_, err = tx.Exec(createGroupParticipantsQuery, groupId, participantId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return participantId, tx.Commit()
}

func (r *ParticipantPostgres) Delete(groupId, paricipantId int) error {
	query := fmt.Sprintf("DELETE FROM %s gpl WHERE gpl.group_id = $1 AND gpl.participant_id = $2", groupsParticipantsListsTable)

	_, err := r.db.Exec(query, groupId, paricipantId)
	return err
}

func (r *ParticipantPostgres) GetRecipientInfo(groupId, participantId int) (testApi.RecipientDTO, error) {
	var participant testApi.Participant
	query := fmt.Sprintf("SELECT * FROM %s p WHERE p.id = $1", participantsTable)
	err := r.db.Get(&participant, query, participantId)

	recipientId := participant.RecipientId
	var recipient testApi.Participant
	query = fmt.Sprintf("SELECT * FROM %s p WHERE p.id = $1", participantsTable)
	err = r.db.Get(&recipient, query, recipientId)

	recipientDTO := testApi.RecipientDTO{
		Id:   recipient.Id,
		Name: recipient.Name,
		Wish: recipient.Wish,
	}

	return recipientDTO, err
}
