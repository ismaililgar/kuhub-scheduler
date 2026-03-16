CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    cron_expr TEXT NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPZ DEFAULT now()
)

CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID REFERENCES jobs(id),
    status TEXT NOT NULL,
    started_at TIMESTAMPZ,
    finished_at TIMESTAMPZ,
    error_msg TEXT,
    duration_ms BIGINT
)