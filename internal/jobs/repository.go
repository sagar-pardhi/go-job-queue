package jobs

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(job *Job) error {
	payload, _ := json.Marshal(job.Payload)

	_, err := r.db.Exec(
		context.Background(),
		`
		INSERT INTO Jobs(
			id,
			type,
			payload,
			status,
			retries,
			max_retries
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		job.ID,
		job.Type,
		payload,
		job.Status,
		job.Retries,
		job.MaxRetries,
	)

	return err
}

func (r *Repository) GetByID(id string) (*Job, error) {
	var job Job

	err := r.db.QueryRow(
		context.Background(),
		`
		SELECT id, type, payload, status, retries, max_retries, COALESCE(error, '')
		FROM jobs
		WHERE id = $1
		`,
		id,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.Retries,
		&job.MaxRetries,
		&job.Error,
	)

	if err != nil {
		return nil, err
	}

	return &job, err
}

func (r *Repository) UpdateStatus(id string, status string) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE jobs
		SET Status = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		status,
		id,
	)

	return err
}

func (r *Repository) List() ([]Job, error) {
	rows, err := r.db.Query(
		context.Background(),
		`
		SELECT id, type, payload, status, retries, max_retries, COALESCE(error, '')
		FROM jobs
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job

		rows.Scan(
			&job.ID,
			&job.Type,
			&job.Payload,
			&job.Status,
			&job.Retries,
			&job.MaxRetries,
			&job.Error,
		)

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *Repository) UpdateFailure(id string, retries int, errMsg string) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE jobs
		SET retries = $1
			error = $2
		WHERE id = $3
		`,
		retries,
		errMsg,
		id,
	)
	return err
}

func (r *Repository) ClearError(id string) error {
	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE jobs
		SET error = NULL
		WHERE id = $1
		`,
		id,
	)
	return err
}
