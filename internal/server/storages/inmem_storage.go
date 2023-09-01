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
	switch mc.Mtype {
	case "gauge":
		m.metrics[mc.ID] = mc
	case "counter":
		if existingMetric, ok := m.metrics[mc.ID]; ok {
			*existingMetric.Delta += *mc.Delta
		} else {
			m.metrics[mc.ID] = mc
		}
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
