package services

import (
	"errors"
	"fmt"
	"strconv"

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

	s.storage.Set(metric.GaugeType.ToString(), name, value)

	return nil
}

func (s *gaugeService) GetValues() (map[string]string, error) {
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

func (s *gaugeService) GetValue(name string) (strValue string, ok bool, err error) {
	value, ok := s.storage.Get(metric.GaugeType.ToString(), name)

	if !ok {
		return "", ok, errors.New("metric " + name + " not found")
	}

	strValue = fmt.Sprintf("%g", value.(float64))

	return
}
