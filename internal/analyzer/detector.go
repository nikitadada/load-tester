package analyzer

import (
	"sort"
	"time"

	"github.com/nikitadada/load-tester/internal/metrics"
)

type DegradationEvent struct {
	Time        time.Time
	BaselineP95 time.Duration
	CurrentP95  time.Duration
	RPS         int
}

type Detector struct {
	cfg Config

	baselineSamples []time.Duration
	baselineReady   bool
	baselineP95     time.Duration

	violationCount int
	triggered      bool
}

func New(cfg Config) *Detector {
	return &Detector{
		cfg: cfg,
	}
}

func (d *Detector) Analyze(
	now time.Time,
	stats metrics.WindowStats,
) *DegradationEvent {

	// если уже нашли деградацию — больше не ищем
	if d.triggered {
		return nil
	}

	// игнорируем пустые окна
	if stats.RPS < d.cfg.MinSamplesPerWindow {
		return nil
	}

	// ===== ЭТАП 1 — сбор baseline =====

	if !d.baselineReady {
		d.baselineSamples = append(d.baselineSamples, stats.P95)

		if len(d.baselineSamples) >= d.cfg.BaselineWindowCount {
			d.baselineP95 = percentileDuration(d.baselineSamples, 0.95)
			d.baselineReady = true
		}

		return nil
	}

	// ===== ЭТАП 2 — проверка нарушения =====

	threshold := time.Duration(
		float64(d.baselineP95) * d.cfg.LatencyFactor,
	)

	if stats.P95 > threshold {
		d.violationCount++
	} else {
		d.violationCount = 0
	}

	// ===== ЭТАП 3 — подтверждение деградации =====

	if d.violationCount >= d.cfg.ViolationWindows {
		d.triggered = true

		return &DegradationEvent{
			Time:        now,
			BaselineP95: d.baselineP95,
			CurrentP95:  stats.P95,
			RPS:         stats.RPS,
		}
	}

	return nil
}

func percentileDuration(values []time.Duration, p float64) time.Duration {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(values))
	copy(sorted, values)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(float64(len(sorted)-1) * p)
	return sorted[index]
}
