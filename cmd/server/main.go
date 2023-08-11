package main

import (
	"errors"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	parseFlags()
	storage := storages.NewMemStorage()
	contr := controllers.NewMetricController(storage)
	router := chi.NewRouter()
	router.Mount("/", contr.Route())

	server := &http.Server{
		Addr:    flagRunAddr,
		Handler: router,
	}

	log.Printf("Серевер запущен на http://%s\n", flagRunAddr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	return nil
}
