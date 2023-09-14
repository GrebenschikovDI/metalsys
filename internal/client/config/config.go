package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type AgentConfig struct {
	ServerAddress  string
	ReportInterval time.Duration
	PollInterval   time.Duration
	HashKey        string
}

var (
	flagSendAddr string
	flagRepInt   string
	flagPollInt  string
	flagKey      string
)

func LoadConfig() (*AgentConfig, error) {
	flag.StringVar(&flagSendAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagRepInt, "r", "10", "interval to send metrics")
	flag.StringVar(&flagPollInt, "p", "2", "interval to update metrics")
	flag.StringVar(&flagKey, "k", "", "sign key")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envKey := os.Getenv("KEY"); envKey != "" {
		flagKey = envKey
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagSendAddr = envRunAddr
	}
	if envRepInt := os.Getenv("REPORT_INTERVAL"); envRepInt != "" {
		flagRepInt = envRepInt
	}
	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		flagPollInt = envPollInt
	}

	pollInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagPollInt))
	if err != nil {
		return nil, fmt.Errorf("cannot parse pollInterval to Duration: %w", err)
	}
	reportInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagRepInt))
	if err != nil {
		return nil, fmt.Errorf("cannot parse reportInterval to Duration: %w", err)
	}

	server := fmt.Sprintf("http://%s/", flagSendAddr)

	cfg := &AgentConfig{
		ServerAddress:  server,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
		HashKey:        flagKey,
	}
	return cfg, nil
}
