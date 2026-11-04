package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type Repo struct {
	*sql.DB
}

func New(db *sql.DB) ProjectRepo {
	return newTraced(&Repo{db})
}

func (r *Repo) Store(ctx context.Context, project *entity.Project) error {
	langs, err := json.Marshal(project.Languages)
	if err != nil {
		return fmt.Errorf("ProjectRepo - Store - json.Marshal: %w", err)
	}

	_, err = r.ExecContext(ctx,
		`INSERT INTO projects (id, project_url, title, description, languages, created_at, visible, custom, body)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		project.ID,
		project.ProjectUrl,
		project.Title,
		project.Description,
		string(langs),
		project.CreatedAt.Format(time.RFC3339),
		project.Visible,
		project.Custom,
		project.Body,
	)
	if err != nil {
		return fmt.Errorf("ProjectRepo - Store - r.ExecContext: %w", err)
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, projectID string) (entity.Project, error) {
	var (
		result    entity.Project
		langs     string
		createdAt string
	)

	err := r.QueryRowContext(ctx,
		`SELECT id, project_url, title, description, languages, created_at, visible, custom, body
		FROM projects
		WHERE id = ?`,
		projectID,
	).Scan(
		&result.ID,
		&result.ProjectUrl,
		&result.Title,
		&result.Description,
		&langs,
		&createdAt,
		&result.Visible,
		&result.Custom,
		&result.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Project{}, entity.ErrProjectNotFound
		}

		return entity.Project{}, fmt.Errorf("ProjectRepo - GetByID - r.QueryRowContext: %w", err)
	}

	err = json.Unmarshal([]byte(langs), &result.Languages)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectRepo - GetByID - json.Unmarshal: %w", err)
	}

	result.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectRepo - GetByID - time.Parse: %w", err)
	}

	return result, nil
}

func (r *Repo) List(ctx context.Context, filter ProjectFilter) ([]entity.Project, int, error) {
	var (
		results []entity.Project
		total   int
	)

	where := "WHERE 1=1"

	if !filter.All {
		where += " AND visible = true"
	}

	err := r.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM projects"+where,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo - List - r.QueryRowContext: %w", err)
	}

	rows, err := r.QueryContext(ctx,
		`SELECT id, project_url, title, description, languages, created_at, visible, custom, body
		FROM projects`+where,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo - List - r.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			result    entity.Project
			langs     string
			createdAt string
		)

		if err := rows.Scan(
			&result.ID,
			&result.ProjectUrl,
			&result.Title,
			&result.Description,
			&langs,
			&createdAt,
			&result.Visible,
			&result.Custom,
			&result.Body,
		); err != nil {
			return nil, 0, fmt.Errorf("ProjectRepo - List - for rows.Scan: %w", err)
		}

		err := json.Unmarshal([]byte(langs), &result.Languages)
		if err != nil {
			return nil, 0, fmt.Errorf("ProjectRepo - List - json.Unmarshal: %w", err)
		}

		result.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("ProjectRepo - List - time.Parse: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo - List - for rows.Next: %w", err)
	}

	return results, total, nil
}

func (r *Repo) Update(ctx context.Context, project *entity.Project) error {
	langs, err := json.Marshal(project.Languages)
	if err != nil {
		return fmt.Errorf("ProjectRepo - Update - for json.Marshal: %w", err)
	}

	result, err := r.ExecContext(ctx,
		`UPDATE projects
			SET project_url = ?, title = ?, description = ?, languages = ?, created_at = ?, visible = ?, custom = ?, body = ?
		 WHERE id = ?`,
		project.ProjectUrl,
		project.Title,
		project.Description,
		string(langs),
		project.CreatedAt.Format(time.RFC3339),
		project.Visible,
		project.Custom,
		project.Body,
		project.ID,
	)
	if err != nil {
		return fmt.Errorf("ProjectRepo - Update - for r.ExecContext: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProjectRepo - Update - for result.RowsAffected: %w", err)
	}

	if n == 0 {
		return entity.ErrProjectNotFound
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, projectID string) error {
	result, err := r.ExecContext(ctx,
		`DELETE FROM projects
		WHERE id = ?`,
		projectID,
	)
	if err != nil {
		return fmt.Errorf("ProjectRepo - Delete - r.ExecContext: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ProjectRepo - Delete - result.RowsAffected: %w", err)
	}

	if n == 0 {
		return entity.ErrProjectNotFound
	}

	return nil
}

func (r *Repo) UpStore(ctx context.Context, project *entity.Project) error {
	langs, err := json.Marshal(project.Languages)
	if err != nil {
		return fmt.Errorf("ProjectRepo - UpStore - json.Marshal: %w", err)
	}

	_, err = r.ExecContext(ctx,
		`
		INSERT INTO projects (id, project_url, title, description, languages, created_at, visible, custom, body)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (project_url) DO UPDATE SET
		title = excluded.title,
		description = excluded.description,
		languages = excluded.languages,
		created_at = excluded.created_at,
		visible = excluded.visible,
		custom = excluded.custom,
		body = excluded.body
		WHERE custom = 0
		`,
		project.ID,
		project.ProjectUrl,
		project.Title,
		project.Description,
		string(langs),
		project.CreatedAt.Format(time.RFC3339),
		project.Visible,
		project.Custom,
		project.Body,
	)
	if err != nil {
		return fmt.Errorf("ProjectRepo - UpStore - r.ExecContext: %w", err)
	}

	return nil
}
