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
	"github.com/ismaililgar/kuhub-scheduler/internal/store"
	"github.com/ismaililgar/kuhub-scheduler/internal/worker"
	"github.com/robfig/cron/v3"
)

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

	// Store bağlantısı
	connStr := "postgres://admin:123456@localhost:5432/kuhub_scheduler?sslmode=disable"
	st, err := store.New(ctx, connStr)
	if err != nil {
		slog.Error("DB bağlantısı kurulamadı", "error", err)
		os.Exit(1)
	}
	defer st.Close()
	slog.Info("DB bağlantısı kuruldu")

	// Registry
	registry := job.NewRegistry()
	registry.Register(&EmailReportJob{})
	registry.Register(&CleanupJob{})

	// Worker pool
	pool := worker.NewPool(3, 10)
	pool.Start(ctx)

	// Cron scheduler
	c := cron.New(cron.WithSeconds())

	for _, j := range registry.All() {
		j := j
		c.AddFunc("*/1 * * * * *", func() {
			start := time.Now()
			slog.Info("Cron tetikledi", "job", j.Name())

			jobCtx, jobCancel := context.WithTimeout(ctx, 30*time.Second)
			defer jobCancel()

			// 1. running olarak kaydet
			if err := st.SaveExecution(jobCtx, store.Execution{
				JobName:   j.Name(),
				Status:    "running",
				StartedAt: start,
			}); err != nil {
				slog.Error("Execution kaydedilemedi", "job", j.Name(), "error", err)
				return
			}

			// 2. job'ı çalıştır — bitene kadar bekler
			execErr := pool.Submit(j)

			// 3. sonuca göre güncelle
			status := "success"
			errMsg := ""
			if execErr != nil {
				status = "failed"
				errMsg = execErr.Error()
			}

			if err := st.UpdateExecution(jobCtx, store.Execution{
				JobName:    j.Name(),
				Status:     status,
				FinishedAt: time.Now(),
				ErrorMsg:   errMsg,
				DurationMs: time.Since(start).Milliseconds(),
			}); err != nil {
				slog.Error("Execution güncellenemedi", "job", j.Name(), "error", err)
			}
		})
	}

	c.Start()
	slog.Info("Cron scheduler başladı", "jobCount", len(registry.All()))

	<-ctx.Done()
	slog.Info("Kapatma sinyali alındı")

	c.Stop()
	pool.Stop()

	slog.Info("Scheduler gracefully kapatıldı")
}
