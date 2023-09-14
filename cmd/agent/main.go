package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/client/core"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
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
	parseFlags()
	pollInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagPollInt))
	if err != nil {
		fmt.Println("Ошибка при парсинге длительности:", err)
		return
	}
	reportInterval, err := time.ParseDuration(fmt.Sprintf("%ss", flagRepInt))
	if err != nil {
		fmt.Println("Ошибка при парсинге длительности:", err)
		return
	}

	server := fmt.Sprintf("http://%s/", flagSendAddr)
	storage := make(map[string]models.Metric)
	var counter int64

	go func() {
		for {
			core.GetRuntimeMetrics(metricNames, storage)
			counter += 1
			storage["PollCount"] = core.GetPollCount(counter)
			storage["RandomValue"] = core.GetRandomValue()
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			controllers.Send(storage, server)
			//controllers.SendJSON(storage, server)
			//controllers.SendSlice(storage, server)
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
