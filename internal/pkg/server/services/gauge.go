package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
)

type gaugeService struct {
	storage memstorage.MemStorage
}

func (s *gaugeService) Update(name string, valStr string) error {
	if len(name) == 0 {
		return errors.New("metric name must be not empty")
	}

	value, err := strconv.ParseFloat(valStr, 64)

	if err != nil {
		return errors.New("incorrect value, must be float")
	}

	valBytes, err := encoding.Float64ToByte(value)
	if err != nil {
		return err
	}

	s.storage.Set(string(metric.GaugeType), name, valBytes)

	return nil
}

func (s *gaugeService) GetValues() (map[string]string, error) {
	valuesBytes, ok := s.storage.GetValuesByNamespace(string(metric.GaugeType))

	values := make(map[string]string)

	if !ok {
		return values, nil
	}

	for k, valueBytes := range valuesBytes {
		value := encoding.Float64FromBytes(valueBytes)

		values[k] = fmt.Sprintf("%g", value)
	}

	return values, nil
}

func (s *gaugeService) GetValue(name string) (value string, ok bool, err error) {
	binaryValue, ok := s.storage.Get(string(metric.GaugeType), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	floatValue := encoding.Float64FromBytes(binaryValue)

	value = fmt.Sprintf("%g", floatValue)

	return
}
