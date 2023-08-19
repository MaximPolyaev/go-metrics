package metricservice

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type MetricService struct {
	storage memStorage
}

type memStorage interface {
	Set(namespace string, key string, val interface{})
	Get(namespace string, key string) (val interface{}, ok bool)
	GetValuesByNamespace(namespace string) (values map[string]interface{}, ok bool)
}

func New(s memStorage) *MetricService {
	return &MetricService{
		storage: s,
	}
}

func (s *MetricService) Update(mm *metric.Metrics) *metric.Metrics {
	switch mm.MType {
	case metric.GaugeType:
		s.storage.Set(mm.MType.ToString(), mm.ID, *mm.Value)
	case metric.CounterType:
		mTypeAsStr := mm.MType.ToString()
		existValue, ok := s.storage.Get(mTypeAsStr, mm.ID)

		if ok {
			*mm.Delta += existValue.(int64)
		}

		s.storage.Set(mTypeAsStr, mm.ID, *mm.Delta)
	}

	return mm
}

func (s *MetricService) Get(mm *metric.Metrics) *metric.Metrics {
	value, ok := s.storage.Get(mm.MType.ToString(), mm.ID)
	if !ok {
		mm.Delta = nil
		mm.Value = nil
		return mm
	}

	switch mm.MType {
	case metric.GaugeType:
		mm.Value = new(float64)
		*mm.Value = value.(float64)
	case metric.CounterType:
		mm.Delta = new(int64)
		*mm.Delta = value.(int64)
	}

	return mm
}

func (s *MetricService) GetValues(mType metric.Type) (map[string]string, error) {
	switch mType {
	case metric.GaugeType:
		return s.getGaugeValues()
	case metric.CounterType:
		return s.getCounterValues()
	}

	return map[string]string(nil), errors.New("unexpected metric type: " + mType.ToString())
}

func (s *MetricService) GetValue(mType metric.Type, name string) (value string, ok bool, err error) {
	switch mType {
	case metric.GaugeType:
		return s.getGaugeValue(name)
	case metric.CounterType:
		return s.getCounterValue(name)
	}

	return "", false, errors.New("unexpected metric type: " + mType.ToString())
}

func (s *MetricService) getGaugeValues() (map[string]string, error) {
	values, ok := s.storage.GetValuesByNamespace(metric.GaugeType.ToString())

	strValues := make(map[string]string)

	if !ok {
		return strValues, nil
	}

	for k, value := range values {
		strValues[k] = fmt.Sprintf("%g", value.(float64))
	}

	return strValues, nil
}

func (s *MetricService) getCounterValues() (map[string]string, error) {
	values, ok := s.storage.GetValuesByNamespace(metric.CounterType.ToString())

	strValues := make(map[string]string)

	if !ok {
		return strValues, nil
	}

	for k, value := range values {
		strValues[k] = strconv.Itoa(int(value.(int64)))
	}

	return strValues, nil
}

func (s *MetricService) getGaugeValue(name string) (strValue string, ok bool, err error) {
	value, ok := s.storage.Get(metric.GaugeType.ToString(), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	strValue = fmt.Sprintf("%g", value.(float64))

	return
}

func (s *MetricService) getCounterValue(name string) (strValue string, ok bool, err error) {
	value, ok := s.storage.Get(metric.CounterType.ToString(), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	strValue = strconv.Itoa(int(value.(int64)))

	return
}
