package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ismaililgar/kuhub-scheduler/internal/job"
	"github.com/ismaililgar/kuhub-scheduler/internal/worker"
)

type EmailReportJob struct{}

func (e *EmailReportJob) Name() string { return "email-report-job" }

func (e *EmailReportJob) Execute(ctx context.Context) error {
	slog.Info("EmailReportJob çalışıyor...")
	select {
	case <-ctx.Done():
		return fmt.Errorf("iptal edildi: %w", ctx.Err())
	default:
		time.Sleep(2 * time.Second)
		slog.Info("EmailReportJob tamamlandı")
		return nil
	}
}

type CleanupJob struct{}

func (c *CleanupJob) Name() string { return "cleanup-job" }

func (c *CleanupJob) Execute(ctx context.Context) error {
	slog.Info("CleanupJob çalışıyor...")

	select {
	case <-ctx.Done():
		return fmt.Errorf("iptal edildi: %w", ctx.Err())
	default:
		time.Sleep(1 * time.Second)
		slog.Info("CleanupJob tamamlandı")
		return nil
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("KUHub Scheduler başlatılıyor...")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	registry := job.NewRegistry()
	registry.Register(&EmailReportJob{})
	registry.Register(&CleanupJob{})

	pool := worker.NewPool(3, 10)
	pool.Start(ctx)

	slog.Info("Worker pool başladı", "workerCount", 3)

	for _, j := range registry.All() {
		slog.Info("Job kuyruğa gönderiliyor", "name", j.Name())
		pool.Submit(j)
	}

	pool.Stop()

	slog.Info("Scheduler has been gracefuly shut down")
}
