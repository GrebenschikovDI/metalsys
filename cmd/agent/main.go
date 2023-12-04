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
		go func() {
			defer wg.Done()
			sendMetricsWorker(storageChan, *cfg)
		}()
	}

	go core.CollectMetrics(storageChan, pollInterval, counter)

	<-ctx.Done()
	close(storageChan)
	wg.Wait()
}

func sendMetricsWorker(ch <-chan map[string]models.Metric, cfg config.AgentConfig) {
	for metrics := range ch {
		controllers.SendSlice(metrics, cfg)
		time.Sleep(cfg.GetReportInterval())
	}
}
