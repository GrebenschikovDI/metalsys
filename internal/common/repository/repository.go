package repository

import (
	"context"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
)

type Repository interface {
	AddMetric(_ context.Context, mc models.Metric) error
	GetMetric(_ context.Context, name string) (value models.Metric, err error)
	GetMetrics(_ context.Context) ([]models.Metric, error)
	AddMetrics(_ context.Context, metrics []models.Metric) error
}
