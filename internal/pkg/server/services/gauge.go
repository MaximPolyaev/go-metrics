package services

import (
	"errors"
	"fmt"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"strconv"
)

type gaugeService struct {
	s memstorage.MemStorage
}

func (updateService *gaugeService) Update(name string, valStr string) error {
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

	updateService.s.Set(string(metric.GaugeType), name, valBytes)

	return nil
}

func (updateService *gaugeService) GetValues() (map[string]string, error) {
	valuesBytes, ok := updateService.s.GetValuesByNamespace(string(metric.GaugeType))

	values := make(map[string]string)

	if !ok {
		return values, nil
	}

	for k, valueBytes := range valuesBytes {
		value := encoding.Float64FromBytes(valueBytes)

		values[k] = fmt.Sprintf("%f", value)
	}

	return values, nil
}
