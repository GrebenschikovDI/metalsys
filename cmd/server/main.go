package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/config"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/storages"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const dirPath = "sql/migrations"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Info("Error loading config", zap.Error(err))
	}
	var storage repository.Repository
	connStr := cfg.GetDsn()
	if connStr == "" {
		storage = storages.NewMemStorage()
	} else {
		db, err := storages.InitDB(context.Background(), connStr, dirPath)
		if err != nil {
			fmt.Println("NO DB")
		}
		storage = db
	}
	filePath := cfg.GetFileStoragePath()
	err = storages.LoadMetrics(cfg.GetRestore(), filePath, storage)
	if err != nil {
		logger.Log.Info("Error reading from file", zap.String("name", filePath))
	}
	interval := cfg.GetStoreInterval()

	go func() {
		for {
			time.Sleep(interval)
			err := storages.SaveMetrics(filePath, storage)
			if err != nil {
				logger.Log.Info("Error writing in file", zap.String("name", filePath))
			}
		}
	}()

	if err := run(storage, *cfg); err != nil {
		panic(err)
	}

	select {}
}

func run(storage repository.Repository, cfg config.ServerConfig) error {
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	ctx := controllers.NewControllerContext(storage, cfg)
	router := controllers.MetricsRouter(ctx)
	address := cfg.GetServerAddress()

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	logger.Log.Info("Running server", zap.String("address", address))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Error within server init", zap.Error(err))
	}

	return nil
}
