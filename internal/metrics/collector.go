package metrics

import (
	"context"
	"sync"
	"time"
)

type Collector struct {
	mu      sync.Mutex
	current *Window

	statsCh chan WindowStats
}

func NewCollector() *Collector {
	return &Collector{
		current: &Window{
			Start:     time.Now(),
			Durations: make([]time.Duration, 0, 1000),
		},
		statsCh: make(chan WindowStats, 100),
	}
}

func (c *Collector) Add(r Result) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.current.Add(r)
}

func (c *Collector) Stats() <-chan WindowStats {
	return c.statsCh
}

func (c *Collector) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			c.rotate()
		}
	}
}

func (c *Collector) rotate() {
	c.mu.Lock()

	old := c.current
	c.current = &Window{
		Start:     time.Now(),
		Durations: make([]time.Duration, 0, 1000),
	}

	c.mu.Unlock()

	stats := old.Stats()

	select {
	case c.statsCh <- stats:
	default:
		// если никто не читает — не блокируемся
	}
}
