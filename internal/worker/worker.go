package worker

import (
	"context"
	"time"

	"github.com/nikitadada/load-tester/internal/grpcclient"
	"github.com/nikitadada/load-tester/internal/metrics"
)

type Worker struct {
	client    grpcclient.Client
	collector *metrics.Collector
}

func New(client grpcclient.Client, collector *metrics.Collector) *Worker {
	return &Worker{
		client:    client,
		collector: collector,
	}
}

func (w *Worker) Do(ctx context.Context) {
	start := time.Now()
	err := w.client.Call(ctx)
	dur := time.Since(start)

	w.collector.Add(metrics.Result{
		Duration: dur,
		Err:      err,
	})
}
