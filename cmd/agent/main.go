package main

import (
	"fmt"
	"reflect"
	"runtime"
)

// server - куда будет отправлен запрос
const server = "http://localhost:8080/"

var metricNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// MetricData - тип данных для сбора и отправки
type MetricData struct {
	Name  string
	Type  string
	Value float64
}

func main() {

	collectMetrics()
}

// сбор метрик

func collectMetrics() []*MetricData {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var metrics []*MetricData
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

		metric := &MetricData{
			Name:  field,
			Type:  "gauge",
			Value: metricValue,
		}
		fmt.Printf("Name: %s, Type: %s, Value: %f\n", metric.Name, metric.Type, metric.Value)
		metrics = append(metrics, metric)
	}
	return metrics
}

// подготовка метрик
// отправка метрик
