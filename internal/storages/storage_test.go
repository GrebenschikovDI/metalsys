package storages

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddGauge(t *testing.T) {
	storage := NewMemStorage()
	storage.AddGauge("metric1", 10.5)

	assert.Equal(t, 10.5, storage.gauges["metric1"], "AddGauge добавление gauge прошло неправильно")
}

func TestAddCounter(t *testing.T) {
	storage := NewMemStorage()
	storage.AddCounter("counter1", 5)

	assert.Equal(t, int64(5), storage.counters["counter1"], "AddCounter сработало неверно")
}

func TestGetMetrics(t *testing.T) {
	storage := NewMemStorage()
	storage.AddGauge("gauge1", 10.5)
	storage.AddCounter("counter1", 5)

	expectedMetrics := []string{
		"/gauge/gauge1/10.500000",
		"/counter/counter1/5",
	}
	metrics := storage.GetMetrics()

	assert.Equal(t, expectedMetrics, metrics, "GetMetrics вернул некорректное значение")
}
