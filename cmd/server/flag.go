package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagStoreInt string
var flagStorePath string
var flagRestore bool
var flagDB string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagStoreInt, "i", "300", "interval to store data")
	flag.StringVar(&flagStorePath, "f", "/tmp/metrics-db.json", "storage path")
	flag.BoolVar(&flagRestore, "r", true, "load saved data from storage")
	flag.StringVar(&flagDB, "d", "", "database address")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

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
			fmt.Println("Error", err)
			return
		}
		flagRestore = boolValue
	}
	if envDataBase := os.Getenv("DATABASE_DSN"); envDataBase != "" {
		flagDB = envDataBase
	} else {
		err := os.Setenv("DATABASE_DSN", flagDB)
		if err != nil {
			return
		}
	}
}
