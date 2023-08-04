package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const serverPort = 8080

func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprint(w, body)
}

func handleMetric(storage controllers.MetricStorage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var metricType string
		var metricName string
		var metricValueStr string

		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		routePattern := regexp.MustCompile(`/update/(?P<type>\w+)/(?P<name>\w+)/(?P<value>[\d.]+)`)
		matches := routePattern.FindStringSubmatch(request.URL.Path)

		if len(matches) == 4 {
			metricType = matches[1]
			metricName = matches[2]
			metricValueStr = matches[3]
		} else {
			http.Error(writer, "Not Fond", http.StatusNotFound)
			return
		}
		if metricName == "" {
			http.Error(writer, "Not Found", http.StatusNotFound)
		}
		switch metricType {
		case "gauge":
			metricValue, err := strconv.ParseFloat(metricValueStr, 64)
			if err != nil {
				http.Error(writer, "Bad Request", http.StatusBadRequest)
				return
			}
			storage.AddGauge(metricName, controllers.Gauge(metricValue))
		case "counter":
			metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				http.Error(writer, "Bad Request", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metricName, controllers.Counter(metricValue))
		default:
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		sendResponse(writer, http.StatusOK, storage.ToString())
	}
}

func main() {
	storage := controllers.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handleMetric(storage))
	port := fmt.Sprintf(":%d", serverPort)
	// Запуск сервера на порту 8080
	log.Printf("Серевер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
