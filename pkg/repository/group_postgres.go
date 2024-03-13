package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"testApi"
)

type GroupPostgres struct {
	db *sqlx.DB
}

func NewGroupPostgres(db *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{db: db}
}

func (r *GroupPostgres) Create(group testApi.Group) (int, error) {
	var id int
	createGroupQuery := fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", groupsTable)

	row := r.db.QueryRow(createGroupQuery, group.Name, group.Description)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *GroupPostgres) GetAll() ([]testApi.Group, error) {
	var groups []testApi.Group
	query := fmt.Sprintf("SELECT * FROM %s", groupsTable)

	err := r.db.Select(&groups, query)
	return groups, err
}

func (r *GroupPostgres) GetById(groupId int) (testApi.GroupDTO, error) {
	var groupDTO testApi.GroupDTO

	var group testApi.Group
	query := fmt.Sprintf(`SELECT * FROM %s g WHERE g.id = $1`, groupsTable)
	err := r.db.Get(&group, query, groupId)

	var participants []testApi.ParticipantDTO

	testQuery := fmt.Sprintf("SELECT p.id, p.name, p.wish, pr.id, pr.name, pr.wish FROM %s p LEFT JOIN %s gpl on p.id = gpl.id LEFT JOIN %s pr on p.recipient_id = pr.id WHERE gpl.group_id = $1", participantsTable, groupsParticipantsListsTable, participantsTable)
	rows, err := r.db.Query(testQuery, groupId)

	for rows.Next() {
		var p testApi.ParticipantDTO
		var recipientId sql.NullInt64
		var recipientName, recipientWish sql.NullString

		if err := rows.Scan(&p.Id, &p.Name, &p.Wish, &recipientId, &recipientName, &recipientWish); err != nil {
			return groupDTO, err
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

	groupDTO = testApi.GroupDTO{
		Group:        group,
		Participants: participants,
	}

	return groupDTO, err
}

func (r *GroupPostgres) Delete(groupId int) error {
	query := fmt.Sprintf("DELETE FROM %s g WHERE g.id = $1", groupsTable)

	_, err := r.db.Exec(query, groupId)
	return err
}

func (r *GroupPostgres) Update(groupId int, input testApi.UpdateGroupInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s g SET %s WHERE g.id = $%d", groupsTable, setQuery, argId)

	args = append(args, groupId)
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
