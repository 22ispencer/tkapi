package db

import (
	"database/sql"
	"strings"

	"github.com/guregu/null/v5"
)

type User struct {
	FirstName          string      `json:"firstName"`
	LastName           null.String `json:"lastName"`
	Badge              null.String `json:"badge"`
	Pin                null.String `json:"pin"`
	FullLegalName      null.String `json:"fullLegalName"`
	Id                 int         `json:"id"`
	LabId              int         `json:"labId"`
	ContactId          null.Int    `json:"contactId"`
	PrimaryContactId   null.Int    `json:"primaryContactId"`
	SecondaryContactId null.Int    `json:"secondaryContactId"`
	ThirdContactId     null.Int    `json:"thirdContactId"`
	TourRoleId         null.Int    `json:"tourRoleId"`
	LabRoleId          null.Int    `json:"labRoleId"`
	IsActive           bool        `json:"isActive"`
}

func (e *Env) GetUsers(labId int, activeOnly bool, labRoleId int) ([]User, error) {
	const query = `
	SELECT u.UserID,
		   u.LabID,
		   u.FirstName,
		   u.LastName,
		   u.ContactID,
		   u.PrimaryContactID,
		   u.SecondContactID,
		   u.ThirdContactID,
		   u.Badge,
		   u.Pin,
		   u.Active,
		   u.TourRoleID,
		   u.FullLegalName,
		   u.LabRoleID
	FROM [User] u
	WHERE (u.LabID = @LabId OR @LabId = 0)
		  AND (u.Active = 1 OR u.Active = @IsActive)
		  AND (u.LabRoleId = @LabRoleId OR @LabRoleId = 0)
	`
	rows, err := e.DB.Query(query, sql.Named("LabId", labId), sql.Named("IsActive", activeOnly), sql.Named("LabRoleId", labRoleId))
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var (
		user  User
		users = []User{}
	)
	for rows.Next() {
		err := rows.Scan(
			&user.Id,
			&user.LabId,
			&user.FirstName,
			&user.LastName,
			&user.ContactId,
			&user.PrimaryContactId,
			&user.SecondaryContactId,
			&user.ThirdContactId,
			&user.Badge,
			&user.Pin,
			&user.IsActive,
			&user.TourRoleId,
			&user.FullLegalName,
			&user.LabRoleId,
		)
		if user.FullLegalName.Valid {
			user.FullLegalName.String = strings.Trim(user.FullLegalName.String, " ")
		}
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (e *Env) GetUserById(userId int) (User, error) {
	const query = `
	SELECT u.UserID,
		   u.LabID,
		   u.FirstName,
		   u.LastName,
		   u.ContactID,
		   u.PrimaryContactID,
		   u.SecondContactID,
		   u.ThirdContactID,
		   u.Badge,
		   u.Pin,
		   u.Active,
		   u.TourRoleID,
		   u.FullLegalName,
		   u.LabRoleID
	FROM [User] u
	WHERE u.UserID = @UserId
	`
	row := e.DB.QueryRow(query, sql.Named("UserId", userId))

	var user User
	err := row.Scan(
		&user.Id,
		&user.LabId,
		&user.FirstName,
		&user.LastName,
		&user.ContactId,
		&user.PrimaryContactId,
		&user.SecondaryContactId,
		&user.ThirdContactId,
		&user.Badge,
		&user.Pin,
		&user.IsActive,
		&user.TourRoleId,
		&user.FullLegalName,
		&user.LabRoleId,
	)
	if err != nil {
		return User{}, err
	}
	if user.FullLegalName.Valid {
		user.FullLegalName.String = strings.Trim(user.FullLegalName.String, " ")
	}
	return user, nil
}
