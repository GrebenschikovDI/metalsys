package storages

import (
	"context"
	"encoding/json"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
	"os"
)

func SaveMetrics(path string, storage repository.Repository) error {
	metrics, err := storage.GetMetrics(context.Background())
	if err != nil {
		return err
	}
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func LoadMetrics(restore bool, path string, storage repository.Repository) error {
	if restore {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var metrics []models.Metric
		if err := json.Unmarshal(data, &metrics); err != nil {
			return err
		}

		for _, value := range metrics {
			storage.AddMetric(context.Background(), value)

		}
	}
	return nil
}
