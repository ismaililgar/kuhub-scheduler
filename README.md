# KUHub Scheduler

A PostgreSQL-backed job scheduler POC written in Go, built as an alternative to Hangfire.

## Motivation

Hangfire is tightly coupled to the .NET runtime, manages its own DB schema, and offers limited
control over its internals. This project solves the same problem using Go's concurrency primitives:

- Goroutine-based worker pool instead of OS threads
- Channel-based job dispatch and result tracking
- PostgreSQL execution history
- Graceful shutdown support

## Architecture

- `cmd/scheduler` — application entry point
- `internal/job` — Job interface and registry
- `internal/worker` — goroutine-based worker pool
- `internal/store` — PostgreSQL execution store

## Getting Started

### Requirements
- Go 1.21+
- PostgreSQL

### Run
```bash
# run migration
psql -U admin -d kuhub_scheduler -f migrations/001_initial.sql

# seed jobs
psql -U admin -d kuhub_scheduler -c "
INSERT INTO jobs (name, cron_expr) VALUES
('cleanup-job', '*/15 * * * * *'),
('email-report-job', '*/15 * * * * *');"

# start
go run ./cmd/scheduler
```

## Why Go over .NET for this?

In .NET, Hangfire uses the ThreadPool — each job occupies an OS thread while waiting on I/O.
Go goroutines start at ~2KB of stack vs ~1MB for OS threads. When a goroutine blocks on I/O
(e.g. writing to Postgres), the Go runtime parks it and reuses the OS thread for another goroutine.
This means 50 concurrent jobs in Go uses a fraction of the resources compared to Hangfire.

## Comparison

| | Hangfire | river | KUHub Scheduler |
|---|---|---|---|
| Runtime | .NET | Go | Go |
| Worker model | ThreadPool | goroutine | goroutine |
| DB schema | auto-managed | auto-managed | explicit migrations |
| Purpose | Production | Production | POC / learning |

> **Note:** Production-ready alternatives like [river](https://github.com/riverqueue/river) and
> [asynq](https://github.com/hibiken/asynq) exist. This project was built to deeply understand
> Go's concurrency model and scheduler mechanics — not to replace them.
