package models

import (
	"studenttaskhub/database"
)

func GetTasksCreatedByUser(username string) ([]Task, error) {

	query := `
	SELECT id, title, description, deadline, priority, status,
	       created_by, claimed_by, created_at, updated_at
	FROM tasks
	WHERE created_by = ?
	`

	rows, err := database.DB.Query(query, username)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.Priority,
			&task.Status,
			&task.CreatedBy,
			&task.ClaimedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
