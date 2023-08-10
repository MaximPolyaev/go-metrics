package services

import (
	"errors"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
)

type counterService struct {
	storage memstorage.MemStorage
}

func (s *counterService) Update(name string, valStr string) error {
	if len(name) == 0 {
		return errors.New("metric name must be not empty")
	}

	value, err := strconv.Atoi(valStr)

	if err != nil {
		return errors.New("incorrect value, must be int")
	}

	sCategory := string(metric.CounterType)

	existValue, ok := s.storage.Get(sCategory, name)

	if ok {
		value += existValue.(int)
	}

	s.storage.Set(sCategory, name, value)

	return nil
}

func (s *counterService) GetValues() (map[string]string, error) {
	values, ok := s.storage.GetValuesByNamespace(string(metric.CounterType))

	strValues := make(map[string]string)

	if !ok {
		return strValues, nil
	}

	for k, value := range values {
		strValues[k] = strconv.Itoa(value.(int))
	}

	return strValues, nil
}

func (s *counterService) GetValue(name string) (strValue string, ok bool, err error) {
	value, ok := s.storage.Get(string(metric.CounterType), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	strValue = strconv.Itoa(value.(int))

	return
}
