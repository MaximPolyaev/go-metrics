package services

import (
	"errors"
	"fmt"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"strconv"
)

type counterService struct {
	s memstorage.MemStorage
}

func (updateService *counterService) Update(name string, valStr string) error {
	if len(name) == 0 {
		return errors.New("metric name must be not empty")
	}

	value, err := strconv.Atoi(valStr)

	if err != nil {
		return errors.New("incorrect value, must be int")
	}

	sCategory := string(metric.CounterType)

	existValueAsBytes, ok := updateService.s.Get(sCategory, name)

	if ok {
		existValue, err := encoding.IntFromBytes(existValueAsBytes)
		if err != nil {
			return err
		}

		value += existValue
	}

	updateService.s.Set(sCategory, name, encoding.IntToByte(value))

	return nil
}

func (updateService *counterService) GetValues() (map[string]string, error) {
	valuesBytes, ok := updateService.s.GetValuesByNamespace(string(metric.CounterType))

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
