package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// AgentConfig представляет конфигурацию агента.
type AgentConfig struct {
	serverAddress  string        // Адрес и порт сервера.
	reportInterval time.Duration // Интервал с которым отправляются данные.
	pollInterval   time.Duration // Интервал с которым собираются данные.
	hashKey        string        // Ключ для подписи данных.
	rateLimit      int           //Количетво одноаременно исходящих запросов на сервер.
}

// Константы с значениями по умолчанию.
const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
	defaultHashKey        = ""
	defaultRateLimit      = 0
)

// LoadConfig загружает конфигурацию агента из флагов командной строки и переменных окружения.
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

// configureFlags настраивает флаги командной строки для конфигурации.
func (c *AgentConfig) configureFlags() error {
	flag.StringVar(&c.hashKey, "k", defaultHashKey, "sign key")
	serverAddress := flag.String("a", defaultServerAddress, "address and port to run server")
	reportInterval := flag.String("r", defaultReportInterval.String(), "interval to send metrics")
	pollInterval := flag.String("p", defaultPollInterval.String(), "interval to update metrics")
	flag.IntVar(&c.rateLimit, "l", defaultRateLimit, "requests limit")
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

// configureEnvVars настраивает конфигурацию из переменных окружения.
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
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		rateLimit, err := strconv.Atoi(envRateLimit)
		if err != nil {
			return err
		}
		c.rateLimit = rateLimit
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

// GetServerAddress возврщает адрес и порт сервера.
func (c *AgentConfig) GetServerAddress() string {
	return c.serverAddress
}

// GetReportInterval возвращает интервал с которым отправляются данные.
func (c *AgentConfig) GetReportInterval() time.Duration {
	return c.reportInterval
}

// GetPollInterval возвращает интервал с которым собираются данные.
func (c *AgentConfig) GetPollInterval() time.Duration {
	return c.pollInterval
}

// GetHashKey возвращает ключ для подписи данных.
func (c *AgentConfig) GetHashKey() string {
	return c.hashKey
}

// GetRateLimit возвращает количество исходящих запросов
func (c *AgentConfig) GetRateLimit() int {
	return c.rateLimit
}
