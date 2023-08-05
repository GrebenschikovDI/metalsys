package main

import (
	"fmt"
	"reflect"
	"runtime"
)

// server - куда будет отправлен запрос
const server = "http://localhost:8080/"

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
	field := "Alloc"
	metrics := []*MetricData{
		{"Alloc", "gauge", float64(reflect.ValueOf(memStats).FieldByName(field).Uint())},
	}
	fmt.Println(*metrics[0])
	return metrics
}

// подготовка метрик
// отправка метрик
