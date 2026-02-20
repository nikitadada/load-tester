package metrics

import "time"

type Result struct {
	Duration time.Duration
	Err      error
}
