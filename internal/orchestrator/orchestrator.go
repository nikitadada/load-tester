package orchestrator

import (
	"context"
	"fmt"
	"sync"

	"github.com/nikitadada/load-tester/internal/config"
	"github.com/nikitadada/load-tester/internal/grpcclient"
	"github.com/nikitadada/load-tester/internal/loadgen"
	"github.com/nikitadada/load-tester/internal/metrics"
	"github.com/nikitadada/load-tester/internal/worker"
)

type Orchestrator struct {
	cfg config.Config
}

func New(cfg config.Config) *Orchestrator {
	return &Orchestrator{cfg: cfg}
}

func (o *Orchestrator) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), o.cfg.Duration)
	defer cancel()

	collector := metrics.NewCollector()
	client := grpcclient.NewDummy()

	workers := make([]*worker.Worker, o.cfg.Workers)
	for i := range workers {
		workers[i] = worker.New(client, collector)
	}

	var wg sync.WaitGroup
	jobCh := make(chan struct{}, o.cfg.Workers)

	// worker pool
	for _, w := range workers {
		wg.Add(1)
		go func(w *worker.Worker) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-jobCh:
					w.Do(ctx)
				}
			}
		}(w)
	}

	scheduler := loadgen.New(o.cfg.RPS)

	go scheduler.Run(ctx, func() {
		select {
		case jobCh <- struct{}{}:
		default:
			// очередь переполнена — drop (важно для стабильного RPS)
		}
	})

	<-ctx.Done()
	wg.Wait()

	total, errors, avg := collector.Summary()

	fmt.Println("==== TEST RESULT ====")
	fmt.Println("total:", total)
	fmt.Println("errors:", errors)
	fmt.Println("avg latency:", avg)
}
