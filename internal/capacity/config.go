package capacity

import "time"

type Config struct {
	MinRPS int
	MaxRPS int

	TestDuration time.Duration

	Precision int           // остановка когда high-low <= precision
	Cooldown  time.Duration // пауза между тестами
}

func DefaultConfig() Config {
	return Config{
		MinRPS:       50,
		MaxRPS:       2000,
		TestDuration: 15 * time.Second,
		Precision:    10,
		Cooldown:     3 * time.Second,
	}
}
