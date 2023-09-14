package controllers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/retry"
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
	url := fmt.Sprintf("%supdate/", server)

	for _, metric := range storage {
		compressedData, err := compressData(metric)
		if err != nil {
			fmt.Println("Ошибка при сжатии данных:", err)
			return
		}

		response, err := sendRequest(url, compressedData)
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

func SendSlice(storage map[string]models.Metric, server string) {
	var metrics []models.Metric

	url := fmt.Sprintf("%supdates/", server)

	for _, value := range storage {
		metrics = append(metrics, value)
	}

	if len(metrics) == 0 {
		return
	}

	compressedData, err := compressData(metrics)
	if err != nil {
		fmt.Println("Ошибка при сжатии данных:", err)
		return
	}

	response, err := SendWithRetry(context.Background(), url, compressedData, 3)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println("Неожиданный ответ:", response.StatusCode)
		return
	}
}

func compressData(data interface{}) ([]byte, error) {
	requestData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var compressedData bytes.Buffer
	compressor := gzip.NewWriter(&compressedData)
	_, err = compressor.Write(requestData)
	if err != nil {
		return nil, err
	}
	compressor.Close()
	return compressedData.Bytes(), nil
}

func sendRequest(url string, requestData []byte) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")
	return client.Do(request)
}

func SendWithRetry(ctx context.Context, url string, requestData []byte, retries int) (*http.Response, error) {
	return retry.Retry(ctx, func() (*http.Response, error) {
		return sendRequest(url, requestData)
	}, retries)
}
