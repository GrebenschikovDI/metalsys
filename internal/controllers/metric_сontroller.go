package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricStorage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetMetrics() []string
	ToString() string
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
	fmt.Println(c.storage.ToString())
	sendResponse(writer, http.StatusOK, c.storage.ToString())
}
