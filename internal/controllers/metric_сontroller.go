package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
	"strconv"
)

type MetricStorage interface {
	AddGauge(ctx context.Context, name string, value float64) error
	AddCounter(ctx context.Context, name string, value int64) error
	//GetMetrics() []string
	GetGauges(_ context.Context) (map[string]float64, error)
	GetCounters(_ context.Context) (map[string]int64, error)
	GetGaugeByName(ctx context.Context, name string) (value float64, err error)
	GetCounterByName(ctx context.Context, name string) (value int64, err error)
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
	r.Use(middleware.Recoverer)
	r.Post("/update/{type}/{name}/{value}", c.handleUpdate)
	r.Get("/", c.mainHandler)
	r.Get("/value/{type}/{name}", c.valueHandler)
	r.Post("/update/", c.jsonUpdateHandler)
	r.Post("/value/", c.jsonValueHandler)
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
		c.storage.AddGauge(request.Context(), metricName, metricValue)
	case "counter":
		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		c.storage.AddCounter(request.Context(), metricName, metricValue)
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
		value, err := c.storage.GetGaugeByName(request.Context(), metricName)
		if err != nil {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		answer = strconv.FormatFloat(value, 'f', -1, 64)
	case "counter":
		value, err := c.storage.GetCounterByName(request.Context(), metricName)
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
	gauges, err := c.storage.GetGauges(request.Context())
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	counters, err := c.storage.GetCounters(request.Context())
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var metricList []string
	for name, value := range gauges {
		metricList = append(metricList, fmt.Sprintf("gauge/%s/%f", name, value))
	}
	for name, value := range counters {
		metricList = append(metricList, fmt.Sprintf("counter/%s/%d", name, value))
	}

	tmpl := template.Must(template.ParseFiles("internal/templates/metricList.html"))
	err = tmpl.Execute(writer, metricList)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *MetricController) jsonUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	var metric models.Metrics
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&metric); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	metricType := metric.Mtype
	switch metricType {
	case "gauge":
		c.storage.AddGauge(request.Context(), metric.ID, *metric.Value)
	case "counter":
		c.storage.AddCounter(request.Context(), metric.ID, *metric.Delta)
	default:
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(metric); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (c *MetricController) jsonValueHandler(writer http.ResponseWriter, request *http.Request) {
	var metric models.Metrics
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&metric); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	metricType := metric.Mtype
	var response models.Metrics
	switch metricType {
	case "gauge":
		value, err := c.storage.GetGaugeByName(request.Context(), metric.ID)
		if err != nil {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		response = models.Metrics{
			ID:    metric.ID,
			Mtype: metricType,
			Delta: nil,
			Value: &value,
		}
	case "counter":
		value, err := c.storage.GetCounterByName(request.Context(), metric.ID)
		if err != nil {
			http.Error(writer, "Not Found", http.StatusNotFound)
			return
		}
		response = models.Metrics{
			ID:    metric.ID,
			Mtype: metricType,
			Delta: &value,
			Value: nil,
		}
	default:
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(response); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
