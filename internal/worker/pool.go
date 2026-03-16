package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ismaililgar/kuhub-scheduler/internal/job"
)

type JobRequest struct {
	Job        job.Job
	ResultChan chan error
}

type Pool struct {
	jobChan chan JobRequest
	workers int
}

func NewPool(workers, bufferSize int) *Pool {
	return &Pool{
		jobChan: make(chan JobRequest, bufferSize),
		workers: workers,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		workerID := i
		go func() {
			slog.Info("Worker başlatıldı", "workerId", workerID)
			for {
				select {
				case req, ok := <-p.jobChan:
					if !ok {
						return
					}
					slog.Info("Worker job alıyor", "workerId", workerID, "job", req.Job.Name())
					err := req.Job.Execute(ctx)
					if err != nil {
						slog.Error("Worker job başarısız", "workerID", workerID, "job", req.Job.Name(), "error", err)
					} else {
						slog.Info("Worker job tamamlandı", "workerID", workerID, "job", req.Job.Name())
					}
					req.ResultChan <- err
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *Pool) Submit(j job.Job) error {
	resultChan := make(chan error, 1)

	req := JobRequest{Job: j, ResultChan: resultChan}
	select {
	case p.jobChan <- req:
		return <-resultChan
	default:
		return fmt.Errorf("job kuyruğu dolu, %s kabul edilemedi", j.Name())
	}
}

func (p *Pool) Stop() {
	close(p.jobChan)
}
