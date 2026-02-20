package main

import (
	"time"

	"github.com/nikitadada/load-tester/internal/config"
	"github.com/nikitadada/load-tester/internal/orchestrator"
)

func main() {
	cfg := config.Config{
		TargetAddr: "localhost:50051",
		RPS:        10,
		Duration:   10 * time.Second,
		Workers:    20,
	}

	orch := orchestrator.New(cfg)
	orch.Run()
}
