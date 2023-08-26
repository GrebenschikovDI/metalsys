package main

import (
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/storages"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	parseFlags()
	storage := storages.NewMemStorage()
	err := storages.LoadMetrics(flagRestore, flagStorePath, storage)
	if err != nil {
		logger.Log.Info("Error reading from file", zap.String("name", flagStorePath))
	}

	storeInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagStoreInt))
	if err != nil {
		logger.Log.Info("Ошибка при парсинге длительности:", zap.Error(err))
	}

	go func() {
		for {
			time.Sleep(storeInterval)
			err := storages.SaveMetrics(flagStorePath, storage)
			if err != nil {
				logger.Log.Info("Error writing in file", zap.String("name", flagStorePath))
			}
		}
	}()

	if err := run(storage); err != nil {
		panic(err)
	}

	select {}
}

func run(storage *storages.MemStorage) error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	ctx := controllers.NewControllerContext(storage)
	router := controllers.MetricsRouter(ctx)

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
