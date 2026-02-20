package metrics

import "time"

type Window struct {
	Start     time.Time
	Durations []time.Duration
	Errors    int
	Count     int
}

func (w *Window) Add(r Result) {
	w.Count++
	if r.Err != nil {
		w.Errors++
	}
	w.Durations = append(w.Durations, r.Duration)
}

type WindowStats struct {
	RPS    int
	Errors int
	P50    time.Duration
	P95    time.Duration
	P99    time.Duration
}

func (w *Window) Stats() WindowStats {
	return WindowStats{
		RPS:    w.Count,
		Errors: w.Errors,
		P50:    Percentile(w.Durations, 0.50),
		P95:    Percentile(w.Durations, 0.95),
		P99:    Percentile(w.Durations, 0.99),
	}
}
