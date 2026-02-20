package metrics

import (
	"sort"
	"time"
)

type Result struct {
	Duration time.Duration
	Err      error
}

func Percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(float64(len(sorted)-1) * p)
	return sorted[index]
}
