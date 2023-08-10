package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
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

	existValueAsBytes, ok := s.storage.Get(sCategory, name)

	if ok {
		existValue, err := encoding.IntFromBytes(existValueAsBytes)
		if err != nil {
			return err
		}

		value += existValue
	}

	s.storage.Set(sCategory, name, encoding.IntToByte(value))

	return nil
}

func (s *counterService) GetValues() (map[string]string, error) {
	valuesBytes, ok := s.storage.GetValuesByNamespace(string(metric.CounterType))

	values := make(map[string]string)

	if !ok {
		return values, nil
	}

	for k, valueBytes := range valuesBytes {
		value, err := encoding.IntFromBytes(valueBytes)

		if err != nil {
			return make(map[string]string), err
		}

		values[k] = fmt.Sprintf("%d", value)
	}

	return values, nil
}

func (s *counterService) GetValue(name string) (value string, ok bool, err error) {
	binaryValue, ok := s.storage.Get(string(metric.CounterType), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	intValue, err := encoding.IntFromBytes(binaryValue)

	if err != nil {
		return
	}

	value = strconv.Itoa(intValue)

	return
}
