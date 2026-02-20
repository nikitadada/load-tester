package capacity

import (
	"fmt"
	"time"
)

type Tester interface {
	RunSingleTest(rps int, duration time.Duration) bool // true если деградация
}

type Searcher struct {
	cfg    Config
	tester Tester
}

func New(cfg Config, tester Tester) *Searcher {
	return &Searcher{
		cfg:    cfg,
		tester: tester,
	}
}

func (s *Searcher) FindCapacity() Result {
	low := s.cfg.MinRPS
	high := s.cfg.MaxRPS

	var lastStable int

	for high-low > s.cfg.Precision {
		mid := (low + high) / 2

		fmt.Println("================================")
		fmt.Println("TESTING RPS =", mid)
		fmt.Println("range:", low, "-", high)

		degraded := s.tester.RunSingleTest(mid, s.cfg.TestDuration)

		if degraded {
			fmt.Println("RESULT: DEGRADED")
			high = mid
		} else {
			fmt.Println("RESULT: STABLE")
			lastStable = mid
			low = mid
		}

		time.Sleep(s.cfg.Cooldown)
	}

	recommended := int(float64(lastStable) * 0.8)

	return Result{
		MaxStableRPS:   lastStable,
		RecommendedRPS: recommended,
	}
}
