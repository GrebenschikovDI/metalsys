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
	//flagDB := Dsn
	var storage repository.Repository
	if cfg.Dsn == "" {
		storage = storages.NewMemStorage()
	} else {
		connStr := cfg.Dsn
		db, err := storages.InitDB(context.Background(), connStr, dirPath)
		if err != nil {
			fmt.Println("NO DB")
		}
		storage = db
	}

	err = storages.LoadMetrics(cfg.Restore, cfg.FileStoragePath, storage)
	if err != nil {
		logger.Log.Info("Error reading from file", zap.String("name", cfg.FileStoragePath))
	}

	go func() {
		for {
			time.Sleep(cfg.StoreInterval)
			err := storages.SaveMetrics(cfg.FileStoragePath, storage)
			if err != nil {
				logger.Log.Info("Error writing in file", zap.String("name", cfg.FileStoragePath))
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

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	logger.Log.Info("Running server", zap.String("address", cfg.ServerAddress))

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Error within server init", zap.Error(err))
	}

	return nil
}
