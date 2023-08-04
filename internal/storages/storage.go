package controllers

import (
	"fmt"
	"strings"
)

type Gauge float64

type Counter int64

// MetricStorage - интерфейс для хранения метрик
type MetricStorage interface {
	AddGauge(name string, value Gauge)
	AddCounter(name string, value Counter)
	ToString() string
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

// ToString - возвращает содержимое storage как строку
func (m *MemStorage) ToString() string {
	var builder strings.Builder

	builder.WriteString("Metrics:\n")
	builder.WriteString("Gauges\n")
	for name, value := range m.gauges {
		builder.WriteString(fmt.Sprintf("%s: %f\n", name, value))
	}
	builder.WriteString("Counters:\n")
	for name, value := range m.counters {
		builder.WriteString(fmt.Sprintf("%s: %d\n", name, value))
	}
	return builder.String()
}
