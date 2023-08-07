package controllers

import (
	"fmt"
	"strings"
)

type MockMetricStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MockMetricStorage) AddGauge(name string, value float64) {
	if m.gauges == nil {
		m.gauges = make(map[string]float64)
	}
	m.gauges[name] = value
}

func (m *MockMetricStorage) AddCounter(name string, value int64) {
	if m.counters == nil {
		m.counters = make(map[string]int64)
	}
	current, ok := m.counters[name]
	if !ok {
		m.counters[name] = value
	} else {
		m.counters[name] = current + value
	}
}

func (m *MockMetricStorage) GetMetrics() []string {
	var results []string
	for name, value := range m.gauges {
		results = append(results, fmt.Sprintf("/gauge/%s/%f", name, value))
	}
	for name, value := range m.counters {
		results = append(results, fmt.Sprintf("/counter/%s/%d", name, value))
	}
	return results
}

func (m *MockMetricStorage) ToString() string {
	var builder strings.Builder

	builder.WriteString("Metrics:\n")
	builder.WriteString("Gauges\n")
	for name, value := range m.gauges {
		builder.WriteString(fmt.Sprintf("%s: %f\n", name, value))
	}
	builder.WriteString("Counters:\n")
	for name, value := range m.counters {
		builder.WriteString(fmt.Sprintf("%s: %d\n", name, value))
	}
	return builder.String()
}

//func TestMetricHandler(t *testing.T) {
//	mockStorage := &MockMetricStorage{}
//	handler := MetricHandler(mockStorage)
//
//	req, err := http.NewRequest(http.MethodPost, "/update/gauge/testMetric/42.0", nil)
//	if err != nil {
//		t.Fatalf("Зпрос не создан: %v", err)
//	}
//
//	rr := httptest.NewRecorder()
//
//	handler(rr, req)
//
//	if rr.Code != http.StatusOK {
//		t.Errorf("Ожидали %v, но получили %v", http.StatusOK, rr.Code)
//	}
//
//	expectedGaugeValue := 42.0
//	if v, ok := mockStorage.gauges["testMetric"]; !ok || v != expectedGaugeValue {
//		t.Errorf("Ожидание %v, но значение не найдено", expectedGaugeValue)
//	}
//}

//func TestMetricSender(t *testing.T) {
//	mockStorage := &MockMetricStorage{}
//	mockStorage.AddGauge("testGauge", 42.0)
//
//	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//	}))
//	defer mockServer.Close()
//
//	MetricSender(mockStorage, mockServer.URL)
//}
