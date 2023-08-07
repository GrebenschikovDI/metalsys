package main

import (
	"flag"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagSendAddr string
var flagRepInt string
var flagPollInt string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagSendAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagRepInt, "r", "10", "interval to send metrics")
	flag.StringVar(&flagPollInt, "p", "2", "interval to update metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
