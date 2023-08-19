package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/core"
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
	//storage := storages.NewMemStorage()

	go func() {
		for {
			//core.UpdateMetrics(metricNames, storage)
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			//controllers.MetricSender(storage, server)
			controllers.JsonMetricUpdate(core.GetJsonMetrics(metricNames), server)
			//fmt.Println(core.GetJsonMetrics(metricNames))
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
