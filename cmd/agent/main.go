package main

import (
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/controllers"
	"github.com/GrebenschikovDI/metalsys.git/internal/core"
	"github.com/GrebenschikovDI/metalsys.git/internal/storages"
	"time"
)

//const server = "http://localhost:8080/"

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
	pollInterval := flagPollInt * time.Second
	reportInterval := flagRepInt * time.Second
	server := fmt.Sprintf("%s/", flagSendAddr)
	storage := storages.NewMemStorage()

	go func() {
		for {
			core.UpdateMetrics(metricNames, storage)
			//fmt.Println(storage.ToString())
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			controllers.MetricSender(storage, server)
			time.Sleep(reportInterval)
		}
	}()

	select {}
}
