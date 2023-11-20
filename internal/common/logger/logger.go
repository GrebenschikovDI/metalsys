package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zapLogger, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zapLogger
	return nil
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		loggingResponseWriter := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   *responseData,
		}

		defer func() {
			duration := time.Since(start)
			Log.Info("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
				zap.Int("status", responseData.status),
				zap.Int("size", responseData.size),
			)
		}()

		next.ServeHTTP(&loggingResponseWriter, r)
	})
}
