package main

import (
	"errors"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}
	storage := storages.NewMemStorage()
	contr := controllers.NewMetricController(storage)
	router := chi.NewRouter()
	router.Use(logger.RequestLogger)
	router.Mount("/", contr.Route())

	server := &http.Server{
		Addr:    flagRunAddr,
		Handler: router,
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Error within server init", zap.Error(err))
	}

	return nil
}
