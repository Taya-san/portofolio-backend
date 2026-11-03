package githubproject

import (
	"context"
	"database/sql"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type Repo struct {
	*sql.DB
}

func New(db *sql.DB) GithubProjectRepo {
	return newTraced(&Repo{db})
}

func (r *Repo) UpStore(ctx context.Context, project *entity.GithubProject) error {
	langs, err := json.Marshal(project.Languages)
	if err != nil {
		return fmt.Errorf("GithubProjectRepo - UpStore - json.Marshal: %w", err)
	}

	_, err = r.ExecContext(ctx,
		`
		INSERT INTO github_projects (id, project_url, title, description, languages, created_at, body)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (project_url) DO UPDATE SET
		title = excluded.title,
		description = excluded.description,
		languages = excluded.languages,
		body = excluded.body
		`,
		project.ID,
		project.ProjectUrl,
		project.Title,
		project.Description,
		string(langs),
		project.CreatedAt.Format(time.RFC3339),
		project.Body,
	)
	if err != nil {
		return fmt.Errorf("GithubProjectRepo - Store - r.ExecContext: %w", err)
	}

	return nil
}

func (r *Repo) GetByProjectUrl(ctx context.Context, projectUrl string) (entity.GithubProject, error) {
	var result entity.GithubProject
	var langs string
	var createdAt string

	err := r.QueryRowContext(ctx,
		`
		SELECT id, project_url, title, description, languages, created_at, body FROM github_projects
		WHERE project_url = ?
		`,
		projectUrl,
	).Scan(
		&result.ID,
		&result.ProjectUrl,
		&result.Title,
		&result.Description,
		&langs,
		&createdAt,
		&result.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.GithubProject{}, entity.ErrProjectNotFound
		}

		return entity.GithubProject{}, fmt.Errorf("GithubProjectRepo - GetByProjectUrl - r.QueryRowContext: %w", err)
	}

	if err := json.Unmarshal([]byte(langs), &result.Languages); err != nil {
		return entity.GithubProject{}, fmt.Errorf("GithubProjectRepo - GetByProjectUrl - json.Unmarshal: %w", err)
	}

	result.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return entity.GithubProject{}, fmt.Errorf("GithubProjectRepo - GetByProjectUrl - time.Parse: %w", err)
	}

	return result, nil
}

func (r *Repo) List(ctx context.Context) ([]entity.GithubProject, int, error) {
	var results []entity.GithubProject
	var total int

	err := r.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM github_projects`,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("GithubProjectRepo - List - r.QueryRowContext: %w", err)
	}

	rows, err := r.QueryContext(ctx,
		`SELECT id, project_url, title, description, languages, created_at, body
		FROM github_projects`,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("GithubProjectRepo - List - r.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result entity.GithubProject
		var langs string
		var createdAt string

		if err := rows.Scan(
			&result.ID,
			&result.ProjectUrl,
			&result.Title,
			&result.Description,
			&langs,
			&createdAt,
			&result.Body,
		); err != nil {
			return nil, 0, fmt.Errorf("GithubProjectRepo - List - for rows.Scan: %w", err)
		}

		err := json.Unmarshal([]byte(langs), &result.Languages)
		if err != nil {
			return nil, 0, fmt.Errorf("GithubProjectRepo - List - json.Unmarshal: %w", err)
		}

		result.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("GithubProjectRepo - List - time.Parse: %w", err)
		}

		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("GithubProjectRepo - List - for rows.Next: %w", err)
	}

	return results, total, nil
}

func (r *Repo) Delete(ctx context.Context, projectID string) error {
	result, err := r.ExecContext(ctx,
		`DELETE FROM github_projects
		WHERE id = ?`,
		projectID,
	)
	if err != nil {
		return fmt.Errorf("GithubProjectRepo - Delete - r.ExecContext: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("GithubProjectRepo - Delete - result.RowsAffected: %w", err)
	}

	if n == 0 {
		return entity.ErrProjectNotFound
	}

	return nil
}
