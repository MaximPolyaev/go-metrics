package metric

import (
	"math/rand"
	"runtime"
)

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
)

type GaugeMap map[string]float64
type CounterMap map[string]int

type Stats struct {
	runtime.MemStats
	PollCount   int
	RandomValue int
}

func (t Type) ToString() string {
	return string(t)
}

func ReadStats(stats *Stats) {
	runtime.ReadMemStats(&stats.MemStats)

	stats.PollCount += 1
	stats.RandomValue = rand.Int()
}

func (s *Stats) GetGaugeMap() GaugeMap {
	return GaugeMap{
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

func (s *Stats) GetCounterMap() CounterMap {
	return CounterMap{
		"PoolCount": s.PollCount,
	}
}
