package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (ts *testStorage) AddGauge(name string, value float64) {
	if ts.gauges == nil {
		ts.gauges = make(map[string]float64)
	}
	ts.gauges[name] = value
}

func (ts *testStorage) AddCounter(name string, value int64) {
	if ts.counters == nil {
		ts.counters = make(map[string]int64)
	}
	current, ok := ts.counters[name]
	if !ok {
		ts.counters[name] = value
	} else {
		ts.counters[name] = current + value
	}

}

func (ts *testStorage) GetMetrics() []string {
	return nil
}

func (ts *testStorage) ToString() string {
	return ""
}

func TestUpdateMetrics(t *testing.T) {

	storage := &testStorage{}

	UpdateMetrics([]string{"HeapAlloc", "HeapSys", "NumGoroutine"}, storage)

	assert.NotNil(t, storage.gauges["HeapAlloc"])
	assert.NotNil(t, storage.gauges["HeapSys"])
	assert.NotNil(t, storage.gauges["NumGoroutine"])
	assert.NotNil(t, storage.counters["PollCount"])
	assert.NotNil(t, storage.gauges["RandomValue"])
}

func TestUpdateMetricsRandomValue(t *testing.T) {

	storage := &testStorage{}

	UpdateMetrics([]string{"HeapAlloc", "RandomValue"}, storage)

	randomValue := storage.gauges["RandomValue"]
	assert.Greater(t, randomValue, float64(0))
	assert.Less(t, randomValue, float64(1))
}

func TestUpdateMetricsPollCount(t *testing.T) {

	storage := &testStorage{}

	for i := 0; i < 5; i++ {
		UpdateMetrics([]string{"PollCount"}, storage)
	}

	assert.Equal(t, int64(5), storage.counters["PollCount"])
}
