package controllers

import (
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
	"github.com/GrebenschikovDI/metalsys.git/internal/server/config"
)

type Metric models.Metric

type ControllerContext struct {
	storage repository.Repository
	cfg     config.ServerConfig
}

func NewControllerContext(storage repository.Repository, cfg config.ServerConfig) *ControllerContext {
	return &ControllerContext{
		storage: storage,
		cfg:     cfg,
	}
}
