package metricservice

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/server/metric"
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

func (s *MetricService) Update(mType metric.Type, name string, value string) error {
	switch mType {
	case metric.GaugeType:
		return s.gaugeUpdate(name, value)
	case metric.CounterType:
		return s.counterUpdate(name, value)
	}

	return errors.New("unexpected metric type: " + mType.ToString())
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

func (s *MetricService) gaugeUpdate(name string, valueStr string) error {
	if len(name) == 0 {
		return errors.New("metric name must be not empty")
	}

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		return errors.New("incorrect value, must be float")
	}

	s.storage.Set(metric.GaugeType.ToString(), name, value)

	return nil
}

func (s *MetricService) counterUpdate(name string, valueStr string) error {
	if len(name) == 0 {
		return errors.New("metric name must be not empty")
	}

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		return errors.New("incorrect value, must be int")
	}

	sCategory := metric.CounterType.ToString()

	existValue, ok := s.storage.Get(sCategory, name)

	if ok {
		value += existValue.(int)
	}

	s.storage.Set(sCategory, name, value)

	return nil
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
		strValues[k] = strconv.Itoa(value.(int))
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

	strValue = strconv.Itoa(value.(int))

	return
}
