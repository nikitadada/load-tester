package loadgen

import (
	"context"
	"time"
)

type Job func()

type Scheduler struct {
	rps int
}

func New(rps int) *Scheduler {
	return &Scheduler{rps: rps}
}

func (s *Scheduler) Run(ctx context.Context, job Job) {
	interval := time.Second / time.Duration(s.rps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job()
		}
	}
}
