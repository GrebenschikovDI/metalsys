package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GrebenschikovDI/metalsys.git/internal/client/config"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/core"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"go.uber.org/zap"
)

var (
	Version = "N/A"
	Date    = "N/A"
	Commit  = "N/A"
)

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", Version)
	fmt.Printf("Build date: %s\n", Date)
	fmt.Printf("Build commit: %s\n", Commit)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Info("Error loading config", zap.Error(err))
	}
	pollInterval := cfg.GetPollInterval()
	rateLimit := cfg.GetRateLimit()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	storageChan := make(chan map[string]models.Metric)
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i <= rateLimit; i++ {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			sendMetricsWorker(ctx, storageChan, *cfg)
		}(ctx)
	}

	go func(ctx context.Context) {
		core.CollectMetrics(ctx, storageChan, pollInterval, counter)
	}(ctx)

	<-ctx.Done()

	wg.Wait()
	close(storageChan)
}

func sendMetricsWorker(ctx context.Context, ch <-chan map[string]models.Metric, cfg config.AgentConfig) {
	for {
		select {
		case metrics, ok := <-ch:
			if !ok {
				return
			}
			controllers.SendSlice(metrics, cfg)
			time.Sleep(cfg.GetReportInterval())
		case <-ctx.Done():
			return
		}

	}
}
