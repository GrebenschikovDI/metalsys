package controllers

import (
	"github.com/GrebenschikovDI/metalsys.git/internal/common/compression"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/hash"
	"github.com/GrebenschikovDI/metalsys.git/internal/common/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricsRouter(ctx *ControllerContext) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(compression.GzipMiddleware)
	r.Use(middleware.Recoverer)
	if ctx.cfg.GetHashKey() != "" {
		r.Use(hash.ValidateHash(ctx.cfg.GetHashKey()))
	}
	r.Get("/", ctx.getRoot)
	r.Get("/value/{type}/{name}", ctx.getValue)
	r.Get("/ping", ctx.ping)
	r.Post("/update/{type}/{name}/{value}", ctx.update)
	r.Post("/update/", ctx.updateJSON)
	r.Post("/value/", ctx.getValueJSON)
	r.Post("/updates/", ctx.updates)
	return r
}
