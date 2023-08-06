package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
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

type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetMetrics() []string
	ToString() string
}

func main() {
	pollInterval := 2 * time.Second
	storage := storages.NewMemStorage()

	go func() {
		for {
			updateMetrics(metricNames, storage)
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			senMetrics(storage)
			time.Sleep(10 * time.Second)
		}
	}()

	select {}
}

// сбор метрик
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

func updateMetrics(metricNames []string, storage MetricStorage) {
	getRuntimeMetrics(metricNames, storage)
	getPollCount(storage)
	getRandomValue(storage)
	fmt.Printf(storage.ToString())
}

func senMetrics(storage MetricStorage) {
	metrics := storage.GetMetrics()
	for _, metric := range metrics {
		url := fmt.Sprintf("%supdate%s", server, metric)
		//fmt.Println(url)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("Ошибка при создании запроса", err)
			return
		}
		request.Header.Set("Content-Type", "text/plain")
		client := &http.Client{Timeout: 5 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		if response.StatusCode != http.StatusOK {
			fmt.Println("Неожиданный ответ:", response.StatusCode)
			return
		}
		response.Body.Close()
	}

}
