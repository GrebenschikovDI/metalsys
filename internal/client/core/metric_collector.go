package core

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"math/rand"
	"reflect"
	"runtime"
)

// GetRuntimeMetrics - собирает метрики из пакета runtime, по списку имен
func GetRuntimeMetrics(metricNames []string, storage map[string]models.Metric) {
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
		metric := models.Metric{
			ID:    field,
			Mtype: "gauge",
			Delta: nil,
			Value: &metricValue,
		}
		storage[field] = metric
	}
}

// GetPollCount - увеличивает PollCount на 1
func GetPollCount(counter int64) models.Metric {
	metric := models.Metric{
		ID:    "PollCount",
		Mtype: "counter",
		Delta: &counter,
		Value: nil,
	}
	return metric
}

// GetRandomValue - добавляет случайное значение
func GetRandomValue() models.Metric {
	randomFloat := rand.Float64()
	metric := models.Metric{
		ID:    "RandomValue",
		Mtype: "gauge",
		Delta: nil,
		Value: &randomFloat,
	}
	return metric
}
