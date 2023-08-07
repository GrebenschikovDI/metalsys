package core

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
)

type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetMetrics() []string
	ToString() string
}

// getRuntimeMetrics - собирает метрики из пакета runtime, по списку имен
func getRuntimeMetrics(metricNames []string, storage MetricStorage) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for _, field := range metricNames {
		value := reflect.ValueOf(memStats).FieldByName(field)
		var metricValue float64

		switch value.Kind() {
		case reflect.Uint64:
			metricValue = float64(value.Uint())
		case reflect.Uint32:
			metricValue = float64(value.Uint())
		case reflect.Float64:
			metricValue = value.Float()
		default:
			fmt.Printf("Неправильный тип метрики %s\n", field)
			continue
		}
		storage.AddGauge(field, metricValue)
	}
}

// getPollCount - увеличивает PollCount на 1
func getPollCount(storage MetricStorage) {
	field := "PollCount"
	storage.AddCounter(field, 1)
}

// getRandomValue - добавляет случайное значение
func getRandomValue(storage MetricStorage) {
	randomFloat := rand.Float64()
	field := "RandomValue"
	storage.AddGauge(field, randomFloat)
}

// UpdateMetrics - обновляет метрики
func UpdateMetrics(metricNames []string, storage MetricStorage) {
	getRuntimeMetrics(metricNames, storage)
	getPollCount(storage)
	getRandomValue(storage)
}
