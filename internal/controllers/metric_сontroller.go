package controllers

import (
	"context"
	"fmt"
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

	t := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Metric List</title>
	</head>
	<body>
		<h1>Metric List</h1>
		<ul>
		{{range .}}
			<li>{{.}}</li>
		{{end}}
		</ul>
	</body>
	</html>
	`
	tmpl, err := template.New("metricList").Parse(t)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(writer, metricList)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
