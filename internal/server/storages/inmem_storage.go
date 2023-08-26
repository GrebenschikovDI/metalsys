package storages

import (
	"context"
	"fmt"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
)

// MemStorage - реализация MetricStorage на основе map
type MemStorage struct {
	metrics map[string]models.Metric
}

// NewMemStorage - создает новое хранлище метрик
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]models.Metric),
	}
}

func (m *MemStorage) AddMetric(_ context.Context, mc models.Metric) error {
	m.metrics[mc.ID] = mc
	switch mc.Mtype {
	case "gauge":
		m.metrics[mc.ID] = mc
	case "counter":
		existingMetric := m.metrics[mc.ID]
		if existingMetric.Delta == nil {
			deltaValue := *mc.Delta
			existingMetric.Delta = &deltaValue
		} else {
			*existingMetric.Delta += *mc.Delta
		}
		m.metrics[mc.ID] = existingMetric
	}
	return nil
}

func (m *MemStorage) GetMetric(_ context.Context, name string) (value models.Metric, err error) {
	value, found := m.metrics[name]
	if !found {
		err = fmt.Errorf("counter with name '%s' not found", name)
	}
	return value, err
}

func (m *MemStorage) GetMetrics(_ context.Context) ([]models.Metric, error) {
	var metrics []models.Metric
	for _, value := range m.metrics {
		metrics = append(metrics, value)
	}
	return metrics, nil
}
