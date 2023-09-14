package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	ServerAddress   string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	Dsn             string
	HashKey         string
}

var (
	flagRunAddr   string
	flagStoreInt  string
	flagStorePath string
	flagRestore   bool
	flagDB        string
	flagKey       string
)

func LoadConfig() (*ServerConfig, error) {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagStoreInt, "i", "300", "interval to store data")
	flag.StringVar(&flagStorePath, "f", "/tmp/metrics-db.json", "storage path")
	flag.BoolVar(&flagRestore, "r", true, "load saved data from storage")
	flag.StringVar(&flagDB, "d", "", "database address")
	flag.StringVar(&flagKey, "k", "", "sign key")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envKey := os.Getenv("KEY"); envKey != "" {
		flagKey = envKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envStoreInt := os.Getenv("STORE_INTERVAL"); envStoreInt != "" {
		flagStoreInt = envStoreInt
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		flagStorePath = envStorePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		boolValue, err := strconv.ParseBool(envRestore)
		if err != nil {
			return nil, fmt.Errorf("cannot parse Restore to bool: %w", err)
		}
		flagRestore = boolValue
	}
	if envDataBase := os.Getenv("DATABASE_DSN"); envDataBase != "" {
		flagDB = envDataBase
	} else {
		os.Setenv("DATABASE_DSN", flagDB)
	}

	storeInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagStoreInt))
	if err != nil {
		return nil, fmt.Errorf("cannot parse StoreInterval to Duration: %w", err)
	}

	cfg := &ServerConfig{
		ServerAddress:   flagRunAddr,
		StoreInterval:   storeInterval,
		FileStoragePath: flagStorePath,
		Restore:         flagRestore,
		Dsn:             flagDB,
		HashKey:         flagKey,
	}

	return cfg, nil
}

func (c *ServerConfig) GetConfig() *ServerConfig {
	return c
}
