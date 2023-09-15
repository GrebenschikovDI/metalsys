package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"html/template"
	"net/http"
	"strconv"
)

func (c *ControllerContext) getRoot(writer http.ResponseWriter, request *http.Request) {
	metrics, err := c.storage.GetMetrics(request.Context())
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var metricList []string
	for _, metric := range metrics {
		switch metric.Mtype {
		case "gauge":
			metricList = append(metricList, fmt.Sprintf("gauge/%s/%f", metric.ID, *metric.Value))
		case "counter":
			metricList = append(metricList, fmt.Sprintf("counter/%s/%d", metric.ID, *metric.Delta))

		}
	}
	writer.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/metricList.html"))
	err = tmpl.Execute(writer, metricList)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (c *ControllerContext) getValue(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")
	var answer string

	metric, err := c.storage.GetMetric(request.Context(), metricName)
	if err != nil {
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	switch metricType {
	case "gauge":
		answer = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case "counter":
		answer = strconv.FormatInt(*metric.Delta, 10)
	default:
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(answer))
}

func (c *ControllerContext) getValueJSON(writer http.ResponseWriter, request *http.Request) {
	var metric models.Metric
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&metric); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var response models.Metric

	value, err := c.storage.GetMetric(request.Context(), metric.ID)
	if err != nil {
		http.Error(writer, "Not Found", http.StatusNotFound)
		return
	}
	response = value

	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(response); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (c *ControllerContext) ping(writer http.ResponseWriter, request *http.Request) {
	db, err := sql.Open("pgx", c.cfg.GetDsn())
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
