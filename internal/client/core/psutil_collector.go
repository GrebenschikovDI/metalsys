package core

import (
	"time"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

const (
	TotalMemoryID    = "TotalMemory"
	FreeMemoryID     = "FreeMemory"
	CPUUtilizationID = "CPUUtilization1"
	GaugeMetricType  = "gauge"
)

func CollectPsutils(metricChan chan<- map[string]models.Metric, interval time.Duration) {
	for {
		storage := make(map[string]models.Metric)
		err := getPsutilsMetrics(storage)
		if err != nil {
			logger.Log.Info("Error collecting metrics", zap.Error(err))
		}
		metricChan <- storage
		time.Sleep(interval)
	}
}

func getPsutilsMetrics(storage map[string]models.Metric) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	totalMemory := float64(memory.Total)
	tm := models.Metric{
		ID:    TotalMemoryID,
		Mtype: GaugeMetricType,
		Delta: nil,
		Value: &totalMemory,
	}
	storage[tm.ID] = tm

	freeMemory := float64(memory.Free)
	fm := models.Metric{
		ID:    FreeMemoryID,
		Mtype: GaugeMetricType,
		Delta: nil,
		Value: &freeMemory,
	}
	storage[fm.ID] = fm

	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	cpUutilization1 := cpuUtilization[0]
	cpuUtil := models.Metric{
		ID:    CPUUtilizationID,
		Mtype: GaugeMetricType,
		Delta: nil,
		Value: &cpUutilization1,
	}
	storage[cpuUtil.ID] = cpuUtil
	return nil
}
