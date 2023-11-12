package gopsutilstats

import (
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type (
	Stats struct {
		TotalMemory     uint64
		FreeMemory      uint64
		CPUutilization1 float64
		log             *logger.Logger
	}

	gaugeMap map[string]float64
)

func New(log *logger.Logger) *Stats {
	return &Stats{log: log}
}

// ReadStats read metric data by mem and cpu stats
func (s *Stats) ReadStats() {
	v, err := mem.VirtualMemory()
	if err != nil {
		s.log.Errorln(err)
		s.TotalMemory = 0
		s.FreeMemory = 0
		s.CPUutilization1 = 0
		return
	}

	s.TotalMemory = v.Total
	s.FreeMemory = v.Free

	percs, err := cpu.Percent(0, false)
	if err != nil {
		s.log.Errorln(err)
		s.TotalMemory = 0
		s.FreeMemory = 0
		s.CPUutilization1 = 0
	}

	for _, p := range percs {
		s.CPUutilization1 = p
	}
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

	return metrics
}

func (s *Stats) getGaugeMap() gaugeMap {
	return gaugeMap{
		"TotalMemory":     float64(s.TotalMemory),
		"FreeMemory":      float64(s.FreeMemory),
		"CPUutilization1": s.CPUutilization1,
	}
}
