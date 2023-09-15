package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	serverAddress  string
	reportInterval time.Duration
	pollInterval   time.Duration
	hashKey        string
}

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
	defaultHashKey        = ""
)

func LoadConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}
	err := cfg.configureFlags()
	if err != nil {
		return nil, err
	}
	err = cfg.configureEnvVars()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *AgentConfig) configureFlags() error {
	flag.StringVar(&c.hashKey, "k", defaultHashKey, "sign key")
	serverAddress := flag.String("a", defaultServerAddress, "address and port to run server")
	reportInterval := flag.String("r", defaultReportInterval.String(), "interval to send metrics")
	pollInterval := flag.String("p", defaultPollInterval.String(), "interval to update metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	c.serverAddress = fmt.Sprintf("http://%s/", *serverAddress)
	duration, err := parseDuration(*reportInterval, defaultReportInterval)
	if err != nil {
		return err
	}
	c.reportInterval = duration
	duration, err = parseDuration(*pollInterval, defaultPollInterval)
	if err != nil {
		return err
	}
	c.pollInterval = duration
	return nil
}

func (c *AgentConfig) configureEnvVars() error {
	if envKey := os.Getenv("KEY"); envKey != "" {
		c.hashKey = envKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.serverAddress = fmt.Sprintf("http://%s/", envRunAddr)
	}
	if envRepInt := os.Getenv("REPORT_INTERVAL"); envRepInt != "" {
		duration, err := parseDuration(envRepInt, defaultReportInterval)
		if err != nil {
			return err
		}
		c.reportInterval = duration
	}
	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		duration, err := parseDuration(envPollInt, defaultReportInterval)
		if err != nil {
			return err
		}
		c.pollInterval = duration
	}
	return nil
}

// parseDuration разбирает строку и возвращает длительность времени.
func parseDuration(value string, defaultValue time.Duration) (time.Duration, error) {
	if _, err := strconv.Atoi(value); err == nil {
		value = fmt.Sprintf("%ss", value)
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue, fmt.Errorf("cannot parse interval to Duration: %w", err)
	}
	return duration, nil
}

func (c *AgentConfig) GetServerAddress() string {
	return c.serverAddress
}
func (c *AgentConfig) GetReportInterval() time.Duration {
	return c.reportInterval
}
func (c *AgentConfig) GetPollInterval() time.Duration {
	return c.pollInterval
}
func (c *AgentConfig) GetHashKey() string {
	return c.hashKey
}
