package memstorage

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_UpdateGaugeMetric(t *testing.T) {
	tests := []struct {
		name        string
		storage     MemStorage
		gauge       metric.Gauge
		wantStorage MemStorage
	}{
		{
			name:    "empty storage",
			storage: MemStorage{gaugeMetrics: make(gaugeMetrics), counterMetrics: make(counterMetrics)},
			gauge:   metric.Gauge{Name: metric.Name("test"), Value: 1},
			wantStorage: MemStorage{
				gaugeMetrics:   gaugeMetrics{"test": metric.Gauge{Name: metric.Name("test"), Value: 1}},
				counterMetrics: make(counterMetrics),
			},
		},
		{
			name: "exist metric in storage",
			storage: MemStorage{
				gaugeMetrics:   gaugeMetrics{"test": metric.Gauge{Name: metric.Name("test"), Value: 0}},
				counterMetrics: make(counterMetrics),
			},
			gauge: metric.Gauge{Name: metric.Name("test"), Value: 1},
			wantStorage: MemStorage{
				gaugeMetrics:   gaugeMetrics{"test": metric.Gauge{Name: metric.Name("test"), Value: 1}},
				counterMetrics: make(counterMetrics),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.UpdateGaugeMetric(tt.gauge)

			assert.Equal(t, tt.storage, tt.wantStorage)
		})
	}
}

func TestMemStorage_UpdateCounterMetric(t *testing.T) {
	tests := []struct {
		name        string
		storage     MemStorage
		counter     metric.Counter
		wantStorage MemStorage
	}{
		{
			name:    "empty storage",
			storage: MemStorage{gaugeMetrics: make(gaugeMetrics), counterMetrics: make(counterMetrics)},
			counter: metric.Counter{Name: metric.Name("test"), Value: 1},
			wantStorage: MemStorage{
				gaugeMetrics:   make(gaugeMetrics),
				counterMetrics: counterMetrics{"test": metric.Counter{Name: metric.Name("test"), Value: 1}},
			},
		},
		{
			name: "exist metric in storage",
			storage: MemStorage{
				gaugeMetrics:   make(gaugeMetrics),
				counterMetrics: counterMetrics{"test": metric.Counter{Name: metric.Name("test"), Value: 1}},
			},
			counter: metric.Counter{Name: metric.Name("test"), Value: 2},
			wantStorage: MemStorage{
				gaugeMetrics:   make(gaugeMetrics),
				counterMetrics: counterMetrics{"test": metric.Counter{Name: metric.Name("test"), Value: 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.UpdateCounterMetric(tt.counter)

			assert.Equal(t, tt.storage, tt.wantStorage)
		})
	}
}
