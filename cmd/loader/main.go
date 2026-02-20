package main

import (
	"github.com/nikitadada/load-tester/internal/config"
	"github.com/nikitadada/load-tester/internal/orchestrator"
)

func main() {
	cfg := config.Config{
		TargetAddr: "localhost:50051",
		Workers:    50,
	}

	orch := orchestrator.New(cfg)

	// обычный тест
	// orch.Run()

	// capacity discovery
	orch.RunCapacitySearch()
}
