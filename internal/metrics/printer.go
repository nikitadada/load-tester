package metrics

import (
	"context"
	"fmt"
)

func StartPrinter(ctx context.Context, statsCh <-chan WindowStats) {
	for {
		select {
		case <-ctx.Done():
			return

		case s := <-statsCh:
			fmt.Printf(
				"RPS=%d errors=%d p50=%v p95=%v p99=%v\n",
				s.RPS,
				s.Errors,
				s.P50,
				s.P95,
				s.P99,
			)
		}
	}
}
