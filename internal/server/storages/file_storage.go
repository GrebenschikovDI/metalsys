package storages

import (
	"context"
	"encoding/json"
	"os"

	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
)

// SaveMetrics writes the metrics stored in the provided repository to a file.
// It fetches all metrics from the repository, marshals them into JSON format,
// and writes the JSON data to the specified file path.
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

// LoadMetrics reads metrics from a specified file and stores them in the repository.
// This function is used to restore metrics from a persistent storage medium.
// It unmarshals the JSON data from the file into metrics and adds them to the repository.
// If the 'restore' flag is false, the function does nothing.
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
