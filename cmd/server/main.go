package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const serverPort = 8080

func main() {
	storage := storages.NewMemStorage()
	contr := controllers.NewMetricController(storage)
	r := chi.NewRouter()
	r.Mount("/", contr.Route())

	port := fmt.Sprintf(":%d", serverPort)

	// Запуск сервера на порту 8080
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	log.Printf("Серевер запущен на http://localhost%s\n", port)
}
