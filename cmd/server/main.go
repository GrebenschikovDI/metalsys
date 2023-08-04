package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"log"
	"net/http"
)

const serverPort = 8080

func main() {
	storage := storages.NewMemStorage()
	handler := controllers.MetricHandler(storage)
	port := fmt.Sprintf(":%d", serverPort)
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", handler)

	// Запуск сервера на порту 8080
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	log.Printf("Серевер запущен на http://localhost%s\n", port)
}
