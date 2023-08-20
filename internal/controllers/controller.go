package controllers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/models"
	"net/http"
	"strconv"
	"time"
)

// sendResponse - добавляет заголовки к ответу
func sendResponse(w http.ResponseWriter, statusCode int, body string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Date", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprint(w, body)
}

func MetricSender(storage MetricStorage, server string) {
	gauges, err := storage.GetGauges(context.TODO())
	if err != nil {
		fmt.Println("Ошибка при запросе gauges", err)
		return
	}
	counters, err := storage.GetCounters(context.TODO())
	if err != nil {
		fmt.Println("Ошибка при запросе counters", err)
		return
	}
	client := &http.Client{Timeout: 10 * time.Second}
	for name, value := range gauges {
		url := fmt.Sprintf("%supdate/gauge/%s/%f", server, name, value)
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
	for name, value := range counters {
		url := fmt.Sprintf("%supdate/counter/%s/%d", server, name, value)
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

func JSONMetricUpdate(storage []models.Metrics, server string) {
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
