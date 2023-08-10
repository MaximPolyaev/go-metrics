package services

import (
	"errors"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
)

type MetricService interface {
	Update(name string, valStr string) error
	GetValues() (map[string]string, error)
	GetValue(name string) (value string, ok bool, err error)
}

func FactoryMetricService(mType metric.Type, s memstorage.MemStorage) (MetricService, error) {
	switch mType {
	case metric.GaugeType:
		return &gaugeService{s: s}, nil
	case metric.CounterType:
		return &counterService{s: s}, nil
	}

	return nil, errors.New("invalid metric type")
}
