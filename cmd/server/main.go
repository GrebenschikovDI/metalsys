package main

import (
	"fmt"
	"log"
	"net/http"
)

const serverPort = 8080

type Gauge float64

type Counter int64

// MetricStorage - интерфейс для хранения метрик
type MetricStorage interface {
	AddGauge(name string, value Gauge)
	AddCounter(name string, value Counter)
}

// MemStorage - реализация MetricStorage на основе map
type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
}

// NewMemStorage - создает новое хранлище метрик
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]Gauge),
		counters: make(map[string]Counter),
	}
}

// AddGauge - добавляет значение типа gauge
func (m *MemStorage) AddGauge(name string, value Gauge) {
	m.gauges[name] = value
}

// AddCounter - добавляет значение типа counter
func (m *MemStorage) AddCounter(name string, value Counter) {
	current, ok := m.counters[name]
	if !ok {
		m.counters[name] = value
	} else {
		m.counters[name] = current + value
	}
}

func main() {

	port := fmt.Sprintf(":%d", serverPort)
	// Запуск сервера на порту 8080
	log.Printf("Серевер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
