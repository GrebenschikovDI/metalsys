package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServerConfig представляет конфигурацию сервера.
type ServerConfig struct {
	serverAddress   string        // Адрес и порт сервера (например, "localhost:8080").
	storeInterval   time.Duration // Интервал сохранения данных.
	fileStoragePath string        // Путь к файловому хранилищу.
	restore         bool          // Флаг, указывающий, следует ли восстанавливать сохраненные данные.
	dsn             string        // Адрес базы данных.
	hashKey         string        // Ключ для подписи данных.
	cryptoKey       string
	trustedSubnet   string
}

type FromFileConfig struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDsn   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	TrustedSubnet string `json:"trusted_subnet"`
}

// Константы с значениями по умолчанию.
const (
	defaultServerAddress   = "localhost:8080"
	defaultInterval        = 300 * time.Second
	defaultRestore         = true
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultDsn             = ""
	defaultHashKey         = ""
	defaultCryptoKey       = ""
	defaultTrustedSubnet   = ""
)

// LoadConfig загружает конфигурацию сервера из флагов командной строки и переменных окружения.
func LoadConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
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
func (c *ServerConfig) configureFlags() error {
	configPath := flag.String("c", "", "config from file")
	// Задаем флаги и их значения по умолчанию.
	flag.StringVar(&c.serverAddress, "a", defaultServerAddress, "address and port to run server")
	flag.StringVar(&c.fileStoragePath, "f", defaultFileStoragePath, "storage path")
	flag.BoolVar(&c.restore, "r", defaultRestore, "load saved data from storage")
	flag.StringVar(&c.dsn, "d", defaultDsn, "database address")
	flag.StringVar(&c.hashKey, "k", defaultHashKey, "sign key")
	storeIntervalStr := flag.String("i", defaultInterval.String(), "interval to store data")
	flag.StringVar(&c.cryptoKey, "crypto-key", defaultCryptoKey, "path to public key")
	flag.StringVar(&c.trustedSubnet, "t", defaultTrustedSubnet, "trusted subnet")
	// Разбираем флаги командной строки.
	flag.Parse()
	if configPath != nil {
		err := c.readConfig(*configPath)
		if err != nil {
			return err
		}
	}
	// Парсим строковое значение интервала и устанавливаем его в StoreInterval.
	duration, err := parseDuration(*storeIntervalStr)
	if err != nil {
		return err
	}
	c.storeInterval = duration
	return nil
}

// configureEnvVars настраивает конфигурацию из переменных окружения.
func (c *ServerConfig) configureEnvVars() error {
	if envKey := os.Getenv("KEY"); envKey != "" {
		c.hashKey = envKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.serverAddress = envRunAddr
	}
	if envStoreInt := os.Getenv("STORE_INTERVAL"); envStoreInt != "" {
		duration, err := parseDuration(envStoreInt)
		if err != nil {
			return err
		}
		c.storeInterval = duration
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		c.fileStoragePath = envStorePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		boolValue, err := parseBool(envRestore)
		if err != nil {
			return err
		}
		c.restore = boolValue
	}
	if envDataBase := os.Getenv("DATABASE_DSN"); envDataBase != "" {
		c.dsn = envDataBase
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		c.cryptoKey = envCryptoKey
	}
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		c.trustedSubnet = envTrustedSubnet
	}
	return nil
}

// parseDuration разбирает строку и возвращает длительность времени.
func parseDuration(value string) (time.Duration, error) {
	if _, err := strconv.Atoi(value); err == nil {
		value = fmt.Sprintf("%ss", value)
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultInterval, fmt.Errorf("cannot parse interval to Duration: %w", err)
	}
	return duration, nil
}

// parseBool разбирает строку и возвращает булево значение.
func parseBool(value string) (bool, error) {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultRestore, fmt.Errorf("cannot parse value to bool: %w", err)
	}
	return boolValue, nil
}

func (c *ServerConfig) readConfig(path string) error {
	if path == "" {
		return nil
	}

	jsonConfig, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	var config FromFileConfig

	err = json.Unmarshal(jsonConfig, &config)

	if err != nil {
		fmt.Println("Error unmarshalling JSON file:", err)
		return err
	}

	if config.Address != "" {
		c.serverAddress = config.Address
	}
	if config.StoreInterval != "" {
		c.storeInterval, err = parseDuration(config.StoreInterval)
		if err != nil {
			return err
		}
	}
	if config.StoreFile != "" {
		c.fileStoragePath = config.StoreFile
	}
	if config.DatabaseDsn != "" {
		c.dsn = config.DatabaseDsn
	}
	if config.CryptoKey != "" {
		c.cryptoKey = config.CryptoKey
	}
	return nil
}

// GetDsn возвращает адрес базы данных.
func (c *ServerConfig) GetDsn() string {
	return c.dsn
}

// GetServerAddress возвращает адрес сервера.
func (c *ServerConfig) GetServerAddress() string {
	return c.serverAddress
}

// GetStoreInterval возвращает интервал сохранения данных.
func (c *ServerConfig) GetStoreInterval() time.Duration {
	return c.storeInterval
}

// GetFileStoragePath возвращает адрес куда сохраняются данные.
func (c *ServerConfig) GetFileStoragePath() string {
	return c.fileStoragePath
}

// GetRestore возвращает флаг, указывающий, следует ли восстанавливать сохраненные данные.
func (c *ServerConfig) GetRestore() bool {
	return c.restore
}

// GetHashKey возвращает ключ для подписи данных.
func (c *ServerConfig) GetHashKey() string {
	return c.hashKey
}

// HasKey возращает true если ключ присутсвует
func (c *ServerConfig) HasKey() bool {
	if c.hashKey == "" {
		return false
	} else {
		return true
	}
}

func (c *ServerConfig) RealIP() string {
	return c.trustedSubnet
}
