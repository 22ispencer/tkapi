package db

import (
	"database/sql"

	"github.com/guregu/null/v5"
)

type Project struct {
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Code        null.String `json:"code"`
	Id          int         `json:"id"`
	LabId       int         `json:"labId"`
	IsActive    bool        `json:"isActive"`
	IsFlagged   bool        `json:"isFlagged"`
}

func (e *Env) GetProjects(labId int, activeOnly bool) ([]Project, error) {
	const query = `
	SELECT p.ProjectID,
		   p.LabID,
		   p.ProjectName,
		   p.Active,
		   p.Description,
		   p.Flagged,
		   p.ProjectCode
	FROM [Project] p
	WHERE (@LabId = p.LabID OR @LabId = 0)
		  AND (p.Active = 1 OR p.Active = @ActiveOnly)
	`
	rows, err := e.DB.Query(query, sql.Named("LabId", labId), sql.Named("ActiveOnly", activeOnly))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		project  Project
		projects = []Project{}
	)
	for rows.Next() {
		err := rows.Scan(
			&project.Id,
			&project.LabId,
			&project.Name,
			&project.IsActive,
			&project.Description,
			&project.IsFlagged,
			&project.Code,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (e *Env) GetProjectById(id int) (Project, error) {
	const query = `
	SELECT p.ProjectID,
		   p.LabID,
		   p.ProjectName,
		   p.Active,
		   p.Description,
		   p.Flagged,
		   p.ProjectCode
	FROM [Project] p
	WHERE p.Projectid = @Id
	`
	row := e.DB.QueryRow(query, sql.Named("Id", id))

	var project Project
	err := row.Scan(
		&project.Id,
		&project.LabId,
		&project.Name,
		&project.IsActive,
		&project.Description,
		&project.IsFlagged,
		&project.Code,
	)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
