package handler_test

import (
	"errors"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type mockMetricService struct{}

func (m *mockMetricService) Update(_ *metric.Metrics) *metric.Metrics {
	return nil
}

func (m *mockMetricService) Get(mm *metric.Metrics) *metric.Metrics {
	return mm
}

func (m *mockMetricService) GetValues(mType metric.Type) (map[string]string, error) {
	switch mType {
	case metric.CounterType:
		return map[string]string{
			"test": "10",
		}, nil
	case metric.GaugeType:
		return map[string]string{
			"test": "1.1",
		}, nil
	}
	return nil, nil
}

func (m *mockMetricService) GetValue(mType metric.Type, name string) (value string, ok bool, err error) {
	if name == "notExist" {
		return "", false, errors.New("")
	}

	switch mType {
	case metric.CounterType:
		if name == "test" {
			return "10", true, nil
		}
	case metric.GaugeType:
		if name == "test" {
			return "1.1", true, nil
		}
	}

	return "", false, nil
}
