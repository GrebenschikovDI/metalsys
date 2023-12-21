//package core
//
//import (
//	"testing"
//	"time"
//
//	"github.com/GrebenschikovDI/metalsys.git/internal/common/models"
//)
//
//func BenchmarkCollectMetrics(b *testing.B) {
//	metricChan := make(chan map[string]models.Metric, 1)
//	interval := 100 * time.Millisecond
//	var counter int64
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		go CollectMetrics(metricChan, interval, counter)
//		<-metricChan
//	}
//}
//
//func TestCollectMetrics(t *testing.T) {
//	metricChan := make(chan map[string]models.Metric, 1)
//	interval := 100 * time.Millisecond
//	var counter int64
//
//	go CollectMetrics(metricChan, interval, counter)
//
//	metrics := <-metricChan
//
//	// Проверяем, присутствуют ли все метрики из списка metricNames
//	for _, name := range metricNames {
//		if _, ok := metrics[name]; !ok {
//			t.Errorf("Метрика %s отсутствует в собранных данных", name)
//		}
//	}
//
//}
