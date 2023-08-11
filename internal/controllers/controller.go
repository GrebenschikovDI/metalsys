package controllers

import (
	"context"
	"fmt"
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
		url := fmt.Sprintf("%supdate/%s/%f", server, name, value)
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
		url := fmt.Sprintf("%supdate/%s/%d", server, name, value)
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
