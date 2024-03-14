package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/rand"
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
	tx, err := r.db.Begin()

	query := fmt.Sprintf("DELETE FROM %s gpl WHERE gpl.group_id = $1 AND gpl.participant_id = $2", groupsParticipantsListsTable)
	_, err = r.db.Exec(query, groupId, paricipantId)

	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE id = $1", participantsTable)
	_, err = r.db.Exec(query, paricipantId)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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

func (r *ParticipantPostgres) Toss(groupId int) ([]testApi.ParticipantDTO, error) {
	tx, err := r.db.Begin()

	var participants []testApi.ParticipantDTO

	var idListOfParticipants []int
	query := fmt.Sprintf("SELECT gpl.participant_id FROM %s gpl WHERE gpl.group_id = $1", groupsParticipantsListsTable)
	err = r.db.Select(&idListOfParticipants, query, groupId)

	if err != nil {
		tx.Rollback()
		return participants, err
	}

	if len(idListOfParticipants) <= 2 {
		tx.Rollback()
		return participants, errors.New("Count of participants less then 2")
	}

	flag, err := GetRecipients(r.db, groupId, idListOfParticipants)
	if err != nil {
		tx.Rollback()
		return participants, err
	}

	if !flag {
		flag, err = GetRecipients(r.db, groupId, idListOfParticipants)
		if err != nil {
			tx.Rollback()
			return participants, err
		}
	}

	testQuery := fmt.Sprintf("SELECT p.id, p.name, p.wish, pr.id, pr.name, pr.wish FROM %s p LEFT JOIN %s gpl on p.id = gpl.id LEFT JOIN %s pr on p.recipient_id = pr.id WHERE gpl.group_id = $1", participantsTable, groupsParticipantsListsTable, participantsTable)
	rows, err := r.db.Query(testQuery, groupId)

	if err != nil {
		tx.Rollback()
		return participants, err
	}

	for rows.Next() {
		var p testApi.ParticipantDTO
		var recipientId sql.NullInt64
		var recipientName, recipientWish sql.NullString

		if err := rows.Scan(&p.Id, &p.Name, &p.Wish, &recipientId, &recipientName, &recipientWish); err != nil {
			tx.Rollback()
			return participants, err
		}

		if int(recipientId.Int64) == 0 {
			p.Recipient = nil
		} else if recipientId.Valid || recipientName.Valid || recipientWish.Valid {
			p.Recipient = &testApi.RecipientDTO{
				Id:   int(recipientId.Int64),
				Name: recipientName.String,
				Wish: recipientWish.String,
			}
		}

		participants = append(participants, p)
	}

	return participants, tx.Commit()
}

func GetRecipients(db *sqlx.DB, groupId int, idListOfParticipants []int) (bool, error) {
	freeParticipants := make([]int, len(idListOfParticipants))
	copy(freeParticipants, idListOfParticipants)

	src := rand.NewSource(rand.Uint64())
	rng := rand.New(src)

	query := fmt.Sprintf("UPDATE %s SET recipient_id = null WHERE id IN (SELECT participant_id FROM %s WHERE group_id = $1)", participantsTable, groupsParticipantsListsTable)
	_, err := db.Exec(query, groupId)

	if err != nil {
		return false, err
	}

	count := 0
	for i := 0; i < len(idListOfParticipants); i++ {
		randomNumber := rng.Intn(len(freeParticipants))
		for idListOfParticipants[i] == freeParticipants[randomNumber] {
			randomNumber = rng.Intn(len(freeParticipants))

			count++
			if count == 10 {
				return false, nil
			}
		}

		query = fmt.Sprintf("UPDATE %s SET recipient_id = $1 WHERE id = $2", participantsTable)
		_, err := db.Exec(query, freeParticipants[randomNumber], idListOfParticipants[i])

		if err != nil {
			return false, err
		}

		freeParticipants = append(freeParticipants[:randomNumber], freeParticipants[randomNumber+1:]...)
	}

	return true, nil
}
