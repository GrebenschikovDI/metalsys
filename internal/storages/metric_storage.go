package storages

import (
	"context"
	"fmt"
	"strings"
)

// MetricStorage - интерфейс для хранения метрик
type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetMetrics() []string
	GetGauge(name string) (value float64, err error)
	GetCounter(ctx context.Context, name string) (value int64, err error)
	ToString() string
}

// MemStorage - реализация MetricStorage на основе map
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage - создает новое хранлище метрик
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// AddGauge - добавляет значение типа gauge
func (m *MemStorage) AddGauge(name string, value float64) {
	m.gauges[name] = value
}

// AddCounter - добавляет значение типа counter
func (m *MemStorage) AddCounter(name string, value int64) {
	current, ok := m.counters[name]
	if !ok {
		m.counters[name] = value
	} else {
		m.counters[name] = current + value
	}
}

// GetMetrics - возращает массив строк, строки в виде /тип/имя/значение
func (m *MemStorage) GetMetrics() []string {
	var results []string
	for name, value := range m.gauges {
		results = append(results, fmt.Sprintf("/gauge/%s/%f", name, value))
	}
	for name, value := range m.counters {
		results = append(results, fmt.Sprintf("/counter/%s/%d", name, value))
	}
	return results
}

func (m *MemStorage) GetGauge(name string) (value float64, err error) {
	value, found := m.gauges[name]
	if !found {
		err = fmt.Errorf("gauge with name '%s' not found", name)
	}
	return value, err
}

func (m *MemStorage) GetCounter(name string) (value int64, err error) {
	value, found := m.counters[name]
	if !found {
		err = fmt.Errorf("counter with name '%s' not found", name)
	}
	return value, err
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
