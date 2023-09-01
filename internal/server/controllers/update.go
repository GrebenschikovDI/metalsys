package controllers

import (
	"encoding/json"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

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
