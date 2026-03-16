package store

import (
	"context"
	"time"
)

type Execution struct {
	JobName    string
	Status     string
	StartedAt  time.Time
	FinishedAt time.Time
	ErrorMsg   string
	DurationMs int64
}

func (s *Store) SaveExecution(ctx context.Context, e Execution) error {
	_, err := s.pool.Exec(ctx, `
							INSERT INTO executions (job_id, status, started_at, finished_at, error_msg, duration_ms)
							SELECT id, $2, $3, $4, $5, $6
							FROM jobs WHERE name = $1`, e.JobName, e.Status, e.StartedAt, e.FinishedAt, e.ErrorMsg, e.DurationMs)
	return err
}

func (s *Store) UpdateExecution(ctx context.Context, e Execution) error {
	_, err := s.pool.Exec(ctx, `
			                 UPDATE executions
							 SET status = $2, finished_at = $3, error_msg = $4, duration_ms = $5
							 WHERE job_id = (SELECT id FROM jobs WHERE name = $1)
							 AND status = 'running'
	`, e.JobName, e.Status, e.FinishedAt, e.ErrorMsg, e.DurationMs)
	return err
}
