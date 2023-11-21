package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/go-chi/chi/v5"
)

// update handles requests to update a single metric value.
// It extracts the metric type, name, and value from URL parameters,
// creates a metric object, and adds or updates it in the storage.
// It responds with appropriate HTTP status codes based on the request processing result.
func (c *ControllerContext) update(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")
	metricValueStr := chi.URLParam(request, "value")

	if metricName == "" {
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	var metric models.Metric
	switch metricType {
	case "gauge":
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		metric = models.Metric{
			ID:    metricName,
			Mtype: metricType,
			Value: &metricValue,
		}
	case "counter":
		metricValue, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		metric = models.Metric{
			ID:    metricName,
			Mtype: metricType,
			Delta: &metricValue,
		}
	default:
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}
	c.storage.AddMetric(request.Context(), metric)
	writer.WriteHeader(http.StatusOK)
}

// updateJSON handles requests with JSON payloads to update a single metric value.
// It decodes the JSON body to a metric object and adds or updates it in the storage.
// It responds with the updated metric in JSON format or appropriate HTTP status codes
// based on the request processing result.
func (c *ControllerContext) updateJSON(writer http.ResponseWriter, request *http.Request) {
	var metric models.Metric
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&metric); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.storage.AddMetric(request.Context(), metric)

	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(metric); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// updates handles batch update requests for multiple metrics through a JSON payload.
// It decodes the JSON body to an array of metric objects and adds or updates them
// in the storage. It responds with an HTTP OK status on successful processing.
func (c *ControllerContext) updates(writer http.ResponseWriter, request *http.Request) {
	var metrics []models.Metric
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&metrics); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.storage.AddMetrics(request.Context(), metrics)
	writer.WriteHeader(http.StatusOK)
}
