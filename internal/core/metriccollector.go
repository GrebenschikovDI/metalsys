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

func getPollCount(storage MetricStorage) {
	field := "PollCount"
	storage.AddCounter(field, 1)
}

func getRandomValue(storage MetricStorage) {
	randomFloat := rand.Float64()
	field := "RandomValue"
	storage.AddGauge(field, randomFloat)
}

func UpdateMetrics(metricNames []string, storage MetricStorage) {
	getRuntimeMetrics(metricNames, storage)
	getPollCount(storage)
	getRandomValue(storage)
}
