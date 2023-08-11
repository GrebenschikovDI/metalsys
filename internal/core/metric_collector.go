package core

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
)

type MetricStorage interface {
	AddGauge(ctx context.Context, name string, value float64) error
	AddCounter(ctx context.Context, name string, value int64) error
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
		err := storage.AddGauge(context.TODO(), field, metricValue)
		if err != nil {
			return
		}
	}
}

// getPollCount - увеличивает PollCount на 1
func getPollCount(storage MetricStorage) {
	field := "PollCount"
	err := storage.AddCounter(context.TODO(), field, 1)
	if err != nil {
		return
	}
}

// getRandomValue - добавляет случайное значение
func getRandomValue(storage MetricStorage) {
	randomFloat := rand.Float64()
	field := "RandomValue"
	err := storage.AddGauge(context.TODO(), field, randomFloat)
	if err != nil {
		return
	}
}

// UpdateMetrics - обновляет метрики
func UpdateMetrics(metricNames []string, storage MetricStorage) {
	getRuntimeMetrics(metricNames, storage)
	getPollCount(storage)
	getRandomValue(storage)
}
