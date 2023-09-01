package controllers

import (
	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/repository"
)

type Metric models.Metric

type ControllerContext struct {
	storage repository.Repository
}

func NewControllerContext(storage repository.Repository) *ControllerContext {
	return &ControllerContext{
		storage: storage,
	}
}
