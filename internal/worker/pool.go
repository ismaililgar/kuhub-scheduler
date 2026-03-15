package worker

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ismaililgar/kuhub-scheduler/internal/job"
)

type Pool struct {
	workerCount int
	jobChan     chan job.Job
	wg          sync.WaitGroup
}

func NewPool(workerCount int, bufferSize int) *Pool {
	return &Pool{
		workerCount: workerCount,
		jobChan:     make(chan job.Job, bufferSize),
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)

		go func(workerId int) {
			defer p.wg.Done()

			slog.Info("Worker başlatıldı", "workerId", workerId)

			for {
				select {
				case j, ok := <-p.jobChan:
					if !ok {
						slog.Info("Worker durdu", "workerID", workerId)
						return
					}
					slog.Info("Worker job alıyor", "workerID", workerId, "job", j.Name())
					if err := j.Execute(ctx); err != nil {
						slog.Error("Job başarısız", "workerID", workerId, "job", j.Name(), "error", err)
					} else {
						slog.Info("Worker job tamamlandı", "workerID", workerId, "job", j.Name())
					}

				case <-ctx.Done():
					slog.Info("Worker context iptal, durdu", "workerID", workerId)
					return
				}
			}
		}(i)
	}
}

func (p *Pool) Submit(j job.Job) {
	p.jobChan <- j
}

func (p *Pool) Stop() {
	close(p.jobChan)
	p.wg.Wait()
	slog.Info("Tüm workerlar durdu")
}
