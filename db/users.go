package db

import "github.com/guregu/null/v5"

type User struct {
	FirstName          string
	LastName           null.String
	Badge              null.String
	Pin                null.String
	FullLegalName      null.String
	Id                 int
	LabId              int
	ContactId          null.Int
	PrimaryContactId   null.Int
	SecondaryContactId null.Int
	ThirdContactId     null.Int
	TourRoleId         null.Int
	LabRoleId          null.Int
	IsActive           bool
}

func (e *Env) Users() ([]User, error) {
	const query = `
	SELECT u.UserID,
		   u.LabID,
		   u.FirstName,
		   u.LastName,
		   u.ContactID,
		   u.PrimaryContactID,
		   u.SecondaryContactID,
		   u.ThirdContactID,
		   u.Badge,
		   u.Pin,
		   u.Active,
		   u.TourRoleID,
		   u.FullLegalName,
		   u.LabRoleID
	FROM [User] u
	`
	rows, err := e.DB.Query(query)
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
		if err != nil {
			return []User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}
