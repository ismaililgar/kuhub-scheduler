package job

import (
	"context"
	"time"
)

type Job interface {
	Name() string
	Execute(ctx context.Context) error
}

type JobDefinition struct {
	ID        string
	Name      string
	CronExpr  string
	Enabled   bool
	CreatedAt time.Time
}

type ExecutionStatus string

const (
	StatusPending ExecutionStatus = "pending"
	StatusRunning ExecutionStatus = "running"
	StatusSuccess ExecutionStatus = "success"
	StatusFailed  ExecutionStatus = "failed"
)

type ExecutionRecord struct {
	ID         string
	JobID      string
	Status     ExecutionStatus
	StartedAt  time.Time
	FinishedAt *time.Time
	Error      *string
	DurationMs int64
}
