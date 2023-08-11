package controllers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetMetrics() []string
	ToString() string
	GetGauge(name string) (value float64, err error)
	GetCounter(name string) (value int64, err error)
}

type MetricController struct {
	storage MetricStorage
}

func NewMetricController(storage MetricStorage) *MetricController {
	return &MetricController{
		storage: storage,
	}
}

func (c *MetricController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", c.handleUpdate)
	r.Get("/", c.mainHandler)
	r.Get("/value/{type}/{name}", c.valueHandler)
	return r
}

func (c *MetricController) handleUpdate(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")
	metricValueStr := chi.URLParam(request, "value")

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
		c.storage.AddGauge(metricName, metricValue)
	case "counter":
		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		c.storage.AddCounter(metricName, metricValue)
	default:
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}
	sendResponse(writer, http.StatusOK, c.storage.ToString())
}

func (c *MetricController) valueHandler(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")
	var answer string
	switch metricType {
	case "gauge":
		value, err := c.storage.GetGauge(metricName)
		if err != nil {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		answer = strconv.FormatFloat(value, 'f', -1, 64)
	case "counter":
		value, err := c.storage.GetCounter(metricName)
		if err != nil {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		answer = strconv.FormatInt(value, 10)
	default:
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	sendResponse(writer, http.StatusOK, answer)
}

func (c *MetricController) mainHandler(writer http.ResponseWriter, request *http.Request) {
	metricList := c.storage.GetMetrics()

	// Генерировать HTML-страницу
	html := "<html><head><title>Metric List</title></head><body><h1>Metric List</h1><ul>"
	for _, metric := range metricList {
		html += "<li>" + metric + "</li>"
	}
	html += "</ul></body></html>"

	// Установить заголовки и отправить ответ
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte(html))
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
