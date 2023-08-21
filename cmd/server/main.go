package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	parseFlags()

	dataChan := make(chan *storages.MemStorage)

	storeInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagStoreInt))
	if err != nil {
		fmt.Println("Ошибка при парсинге длительности:", err)
		return
	}

	go func() {
		if err := run(dataChan); err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			time.Sleep(storeInterval)
			storageToWrite := <-dataChan
			err := os.WriteFile(flagStorePath, storageToWrite.ToJSONBytes(), 0666)
			if err != nil {
				logger.Log.Info("Error writing in file", zap.String("name", flagStorePath))
			}
		}
	}()

	select {}
}

func run(dataChan chan<- *storages.MemStorage) error {
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

	ctx := context.Background()

	if flagRestore {
		data, err := os.ReadFile(flagStorePath)
		if err != nil {
			logger.Log.Info("Error reading from file", zap.String("name", flagStorePath))
		} else {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}

				var jsonValue models.Metrics
				err := json.Unmarshal([]byte(line), &jsonValue)
				if err != nil {
					logger.Log.Info("Error decoding JSON", zap.Error(err))
					continue
				}
				switch jsonValue.Mtype {
				case "gauge":
					storage.AddGauge(ctx, jsonValue.ID, *jsonValue.Value)
				case "counter":
					storage.AddCounter(ctx, jsonValue.ID, *jsonValue.Delta)
				}
			}
		}
	}

	go func() {
		logger.Log.Info("Running server", zap.String("address", flagRunAddr))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("Error within server init", zap.Error(err))
		}
	}()
	go func() {
		for {
			dataChan <- storage
		}
	}()

	return nil
}
