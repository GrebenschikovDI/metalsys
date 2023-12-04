package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// getRoot handles the requests to the root endpoint of the server.
// It retrieves a list of all metrics from the storage and displays them in HTML format.
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

// getValue handles requests to retrieve the value of a specific metric.
// It expects the metric type and name as URL parameters and returns
// the value of the metric in plain text format.
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

// getValueJSON handles JSON formatted requests to retrieve a specific metric.
// It expects a JSON body with metric details, retrieves the metric value,
// and responds with the metric value in JSON format.
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

// ping handles requests to check the database connection.
// It attempts to open a new database connection and pings the database.
// If successful, it returns an HTTP 200 OK status.
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
