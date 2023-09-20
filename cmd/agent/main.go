package main

import (
	"github.com/GrebenschikovDI/metalsys.git/internal/client/config"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/core"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Info("Error loading config", zap.Error(err))
	}
	pollInterval := cfg.GetPollInterval()
	reportInterval := cfg.GetReportInterval()
	rateLimit := cfg.GetRateLimit()

	storageChan := make(chan map[string]models.Metric)
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i <= rateLimit; i++ {
		wg.Add(1)
		go sendMetricsWorker(storageChan, *cfg, &wg, reportInterval)
	}

	go core.CollectMetrics(storageChan, pollInterval, counter)

	go core.CollectPsutils(storageChan, pollInterval)

	wg.Wait()
}

func sendMetricsWorker(ch <-chan map[string]models.Metric, cfg config.AgentConfig, wg *sync.WaitGroup, t time.Duration) {
	defer wg.Done()
	for metrics := range ch {
		controllers.SendSlice(metrics, cfg)
		time.Sleep(t)
	}
}
