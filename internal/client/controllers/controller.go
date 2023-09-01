package controllers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"net/http"
	"time"
)

func Send(metrics map[string]models.Metric, server string) {
	client := &http.Client{Timeout: 10 * time.Second}
	var url string
	for _, metric := range metrics {
		switch metric.Mtype {
		case "gauge":
			url = fmt.Sprintf("%supdate/%s/%s/%f", server, metric.Mtype, metric.ID, *metric.Value)
		case "counter":
			url = fmt.Sprintf("%supdate/%s/%s/%d", server, metric.Mtype, metric.ID, *metric.Delta)
		}
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			fmt.Println("Ошибка при создании запроса", err)
			return
		}
		request.Header.Set("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		if response.StatusCode != http.StatusOK {
			fmt.Println("Неожиданный ответ:", response.StatusCode)
			return
		}
		response.Body.Close()
	}
}

func SendJSON(storage map[string]models.Metric, server string) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%supdate/", server)
	for _, metric := range storage {
		requestData, err := json.Marshal(metric)
		if err != nil {
			fmt.Println("Ошибка при сериализации метрики в JSON:", err)
			return
		}
		var compressedData bytes.Buffer
		compressor := gzip.NewWriter(&compressedData)
		_, err = compressor.Write(requestData)
		if err != nil {
			fmt.Println("Ошибка при сжатии данных:", err)
			return
		}
		compressor.Close()

		request, err := http.NewRequest(http.MethodPost, url, &compressedData)
		if err != nil {
			fmt.Println("Ошибка при создании запроса", err)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			fmt.Println("Неожиданный ответ:", response.StatusCode)
			return
		}
	}
}
