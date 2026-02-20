package metrics

import (
	"sync"
	"time"
)

type Collector struct {
	mu        sync.Mutex
	total     int
	errors    int
	durations []time.Duration
}

func NewCollector() *Collector {
	return &Collector{
		durations: make([]time.Duration, 0, 1000),
	}
}

func (c *Collector) Add(r Result) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.total++
	if r.Err != nil {
		c.errors++
	}
	c.durations = append(c.durations, r.Duration)
}

func (c *Collector) Summary() (total int, errors int, avg time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	total = c.total
	errors = c.errors

	var sum time.Duration
	for _, d := range c.durations {
		sum += d
	}

	if len(c.durations) > 0 {
		avg = sum / time.Duration(len(c.durations))
	}

	return
}
