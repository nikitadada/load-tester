package config

import "time"

type Config struct {
	TargetAddr string        // host:port gRPC сервиса
	RPS        int           // запросов в секунду
	Duration   time.Duration // длительность теста
	Workers    int           // количество воркеров
}
