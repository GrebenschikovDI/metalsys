package main

import (
	"flag"
	"time"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagRepInt time.Duration
var flagPollInt time.Duration

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.DurationVar(&flagRepInt, "r", 10, "interval to send metrics")
	flag.DurationVar(&flagPollInt, "p", 2, "interval to update metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
