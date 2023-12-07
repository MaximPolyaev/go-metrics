// Package defaultstats - это пакет для чтения метрик системы
package defaultstats

import (
	"math/rand"
	"runtime"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type (
	Stats struct {
		runtime.MemStats
		PollCount   int64
		RandomValue int
	}

	gaugeMap   map[string]float64
	counterMap map[string]int64
)

func New() *Stats {
	return &Stats{}
}

// ReadStats - read metric data by runtime stats with update custom fields
func (s *Stats) ReadStats() {
	runtime.ReadMemStats(&s.MemStats)

	s.PollCount += 1
	s.RandomValue = rand.Int()
}

// AsMetrics get stats as metrics
func (s *Stats) AsMetrics() []metric.Metric {
	metrics := make([]metric.Metric, 0, 30)

	for k, v := range s.getGaugeMap() {
		mm := metric.Metric{
			ID:    k,
			MType: metric.GaugeType,
			Value: new(float64),
		}

		*mm.Value = v

		metrics = append(metrics, mm)
	}

	for k, v := range s.getCounterMap() {
		mm := metric.Metric{
			ID:    k,
			MType: metric.CounterType,
			Delta: new(int64),
		}

		*mm.Delta = v

		metrics = append(metrics, mm)
	}

	return metrics
}

func (s *Stats) getGaugeMap() gaugeMap {
	return gaugeMap{
		"Alloc":         float64(s.Alloc),
		"BuckHashSys":   float64(s.BuckHashSys),
		"GCCPUFraction": s.GCCPUFraction,
		"GCSys":         float64(s.GCSys),
		"HeapAlloc":     float64(s.HeapAlloc),
		"HeapIdle":      float64(s.HeapIdle),
		"HeapInuse":     float64(s.HeapInuse),
		"HeapObjects":   float64(s.HeapObjects),
		"HeapReleased":  float64(s.HeapReleased),
		"HeapSys":       float64(s.HeapSys),
		"LastGC":        float64(s.LastGC),
		"Lookups":       float64(s.Lookups),
		"MCacheInuse":   float64(s.MCacheInuse),
		"MCacheSys":     float64(s.MCacheSys),
		"MSpanInuse":    float64(s.MSpanInuse),
		"MSpanSys":      float64(s.MSpanSys),
		"Mallocs":       float64(s.Mallocs),
		"Frees":         float64(s.Frees),
		"NextGC":        float64(s.NextGC),
		"NumForcedGC":   float64(s.NumForcedGC),
		"NumGC":         float64(s.NumGC),
		"OtherSys":      float64(s.OtherSys),
		"PauseTotalNs":  float64(s.PauseTotalNs),
		"StackInuse":    float64(s.StackInuse),
		"StackSys":      float64(s.StackSys),
		"Sys":           float64(s.Sys),
		"TotalAlloc":    float64(s.TotalAlloc),
		"RandomValue":   float64(s.RandomValue),
	}
}

func (s *Stats) getCounterMap() counterMap {
	return counterMap{
		"PollCount": s.PollCount,
	}
}
