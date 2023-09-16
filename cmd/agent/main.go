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

var metricNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Info("Error loading config", zap.Error(err))
	}
	pollInterval := cfg.GetPollInterval()
	reportInterval := cfg.GetReportInterval()

	storage := make(map[string]models.Metric)
	var counter int64
	var m sync.RWMutex

	go func() {
		for {
			m.Lock()
			core.GetRuntimeMetrics(metricNames, storage)
			counter += 1
			storage["PollCount"] = core.GetPollCount(counter)
			storage["RandomValue"] = core.GetRandomValue()
			m.Unlock()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			//controllers.Send(storage, *cfg)
			//controllers.SendJSON(storage, *cfg)
			m.RLock()
			controllers.SendSlice(storage, *cfg)
			m.RUnlock()
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
