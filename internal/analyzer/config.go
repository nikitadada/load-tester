package analyzer

type Config struct {
	BaselineWindowCount int     // сколько окон собираем baseline
	LatencyFactor       float64 // во сколько раз latency должна вырасти
	ViolationWindows    int     // сколько окон подряд должно нарушаться
	MinSamplesPerWindow int     // игнорируем окна с малым RPS
}

func DefaultConfig() Config {
	return Config{
		BaselineWindowCount: 5,
		LatencyFactor:       2.0,
		ViolationWindows:    3,
		MinSamplesPerWindow: 10,
	}
}
