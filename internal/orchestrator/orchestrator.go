package orchestrator

import (
	"context"
	"fmt"
	"github.com/nikitadada/load-tester/internal/analyzer"
	"github.com/nikitadada/load-tester/internal/capacity"
	"sync"
	"time"

	"github.com/nikitadada/load-tester/internal/config"
	"github.com/nikitadada/load-tester/internal/grpcclient"
	"github.com/nikitadada/load-tester/internal/loadgen"
	"github.com/nikitadada/load-tester/internal/metrics"
	"github.com/nikitadada/load-tester/internal/worker"
)

type TestResult struct {
	Degraded bool
}

type Orchestrator struct {
	cfg config.Config
}

func New(cfg config.Config) *Orchestrator {
	return &Orchestrator{cfg: cfg}
}

func (o *Orchestrator) RunCapacitySearch() {
	cfg := capacity.DefaultConfig()
	searcher := capacity.New(cfg, o)

	result := searcher.FindCapacity()

	fmt.Println("================================")
	fmt.Println("CAPACITY RESULT")
	fmt.Println("Max stable RPS:", result.MaxStableRPS)
	fmt.Println("Recommended RPS:", result.RecommendedRPS)
	fmt.Println("================================")
}
func (o *Orchestrator) RunSingleTest(
	rps int,
	duration time.Duration,
) bool {

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	collector := metrics.NewCollector()
	go collector.Start(ctx)

	detector := analyzer.New(analyzer.DefaultConfig())

	client, err := grpcclient.NewPingClient(o.cfg.TargetAddr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	workers := make([]*worker.Worker, o.cfg.Workers)
	for i := range workers {
		workers[i] = worker.New(client, collector)
	}

	jobCh := make(chan struct{}, o.cfg.Workers)

	for _, w := range workers {
		go func(w *worker.Worker) {
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

	scheduler := loadgen.New(rps)

	go scheduler.Run(ctx, func() {
		select {
		case jobCh <- struct{}{}:
		default:
		}
	})

	degraded := false

	for {
		select {
		case <-ctx.Done():
			return degraded

		case s := <-collector.Stats():
			if detector.Analyze(time.Now(), s) != nil {
				degraded = true
			}
		}
	}
}

func (o *Orchestrator) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), o.cfg.Duration)
	defer cancel()

	collector := metrics.NewCollector()
	client, err := grpcclient.NewPingClient(o.cfg.TargetAddr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	detector := analyzer.New(analyzer.DefaultConfig())

	go collector.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case s := <-collector.Stats():

				fmt.Printf(
					"RPS=%d errors=%d p50=%v p95=%v p99=%v\n",
					s.RPS,
					s.Errors,
					s.P50,
					s.P95,
					s.P99,
				)

				event := detector.Analyze(time.Now(), s)
				if event != nil {
					fmt.Println("================================")
					fmt.Println("DEGRADATION DETECTED")
					fmt.Println("baseline p95:", event.BaselineP95)
					fmt.Println("current  p95:", event.CurrentP95)
					fmt.Println("RPS:", event.RPS)
					fmt.Println("time:", event.Time)
					fmt.Println("================================")
				}
			}
		}
	}()

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
}
