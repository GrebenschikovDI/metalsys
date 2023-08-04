package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	ToString() string
}

// sendResponse - добавляет заголовки к ответу
func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprint(w, body)
}

// MetricHandler - обрабатывает POST запрос на сервер
func MetricHandler(storage MetricStorage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var metricType string
		var metricName string
		var metricValueStr string

		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		routePattern := regexp.MustCompile(`/update/(?P<type>\w+)/(?P<name>\w+)/(?P<value>\w+)`)
		matches := routePattern.FindStringSubmatch(request.URL.Path)

		if len(matches) == 4 {
			metricType = matches[1]
			metricName = matches[2]
			metricValueStr = matches[3]
		} else {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		if metricName == "" {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		switch metricType {
		case "gauge":
			metricValue, err := strconv.ParseFloat(metricValueStr, 64)
			if err != nil {
				http.Error(writer, "Bad Request", http.StatusBadRequest)
				return
			}
			storage.AddGauge(metricName, metricValue)
		case "counter":
			metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				http.Error(writer, "Bad Request", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metricName, metricValue)
		default:
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		sendResponse(writer, http.StatusOK, storage.ToString())
	}
}
