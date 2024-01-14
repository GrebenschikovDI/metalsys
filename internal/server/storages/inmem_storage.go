package storages

import (
	"context"
	"fmt"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
)

// MemStorage implements the MetricStorage interface using an in-memory map.
// It is a lightweight storage solution, suitable for temporary storage of metrics.
type MemStorage struct {
	metrics map[string]models.Metric
}

// NewMemStorage creates and returns a new instance of MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]models.Metric),
	}
}

// AddMetric adds a single metric to the MemStorage.
// For 'gauge' type metrics, it replaces the existing metric.
// For 'counter' type metrics, it increments the existing metric by the provided delta.
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

// AddMetrics adds multiple metrics to the MemStorage.
// It iterates over each metric and adds them individually using AddMetric.
func (m *MemStorage) AddMetrics(ctx context.Context, metrics []models.Metric) error {
	for _, metric := range metrics {
		m.AddMetric(ctx, metric)
	}
	return nil
}

// GetMetric retrieves a single metric by its ID from MemStorage.
func (m *MemStorage) GetMetric(_ context.Context, name string) (value models.Metric, err error) {
	value, found := m.metrics[name]
	if !found {
		err = fmt.Errorf("counter with name '%s' not found", name)
	}
	return value, err
}

// GetMetrics retrieves all stored metrics from MemStorage.
func (m *MemStorage) GetMetrics(_ context.Context) ([]models.Metric, error) {
	var metrics []models.Metric
	for _, value := range m.metrics {
		metrics = append(metrics, value)
	}
	return metrics, nil
}
