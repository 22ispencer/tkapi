package db

import (
	"database/sql"

	"github.com/guregu/null/v5"
)

type Task struct {
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Code        null.String `json:"code"`
	Id          int         `json:"id"`
	ProjectId   int         `json:"projectId"`
	IsActive    bool        `json:"isActive"`
}

func (e *Env) GetTasks(projectId int, activeOnly bool) ([]Task, error) {
	const query = `
	SELECT t.TaskID,
		   t.ProjectID,
		   t.TaskName,
		   t.Active,
		   t.Description,
		   t.TaskCode
	FROM [Task] t
	WHERE t.ProjectID = @ProjectId
		  AND (t.Active = 1 OR t.Active = @ActiveOnly)
	`
	rows, err := e.DB.Query(query, sql.Named("ProjectId", projectId), sql.Named("ActiveOnly", activeOnly))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		task  Task
		tasks = []Task{}
	)
	for rows.Next() {
		err := rows.Scan(
			&task.Id,
			&task.ProjectId,
			&task.Name,
			&task.IsActive,
			&task.Description,
			&task.Code,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (e *Env) GetTaskById(taskId int) (Task, error) {
	const query = `
	SELECT t.TaskID,
		   t.ProjectID,
		   t.TaskName,
		   t.Active,
		   t.Description,
		   t.TaskCode
	FROM [Task] t
	WHERE t.TaskID = @TaskId
	`
	row := e.DB.QueryRow(query, sql.Named("TaskId", taskId))

	var task Task
	err := row.Scan(
		&task.Id,
		&task.ProjectId,
		&task.Name,
		&task.IsActive,
		&task.Description,
		&task.Code,
	)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}
