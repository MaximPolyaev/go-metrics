package server

import (
	"errors"
	"strconv"
	"strings"
)

const (
	metricsGaugeType   = "gauge"
	metricsCounterType = "counter"
)

type gaugeMetric struct {
	name  string
	value float64
}

type counterMetric struct {
	name  string
	value int
}

func makeGaugeMetricByURLPath(urlPath string) (*gaugeMetric, error) {
	metricParams, err := urlPathToMetricParamsArr(urlPath)
	if err != nil {
		return nil, err
	}

	name := (*metricParams)[0]
	valueStr := (*metricParams)[1]

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		return nil, errors.New("incorrect value, must be float")
	}

	return &gaugeMetric{
		name:  name,
		value: value,
	}, nil
}

func makeCounterMetricByURLPath(urlPath string) (*counterMetric, error) {
	metricParams, err := urlPathToMetricParamsArr(urlPath)
	if err != nil {
		return nil, err
	}

	name := (*metricParams)[0]
	valueStr := (*metricParams)[1]

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		return nil, errors.New("incorrect value, must be int")
	}

	return &counterMetric{
		name:  name,
		value: value,
	}, nil
}

func urlPathToMetricParamsArr(urlPath string) (*[]string, error) {
	urlPath = strings.Trim(urlPath, "/")

	metricParams := strings.Split(urlPath, "/")

	if len(metricParams) != 2 {
		return nil, errors.New("incorrect metric params count")
	}

	for k, v := range metricParams {
		v = strings.Trim(v, " ")

		if len(v) == 0 {
			return nil, errors.New("is empty metric param")
		}

		metricParams[k] = v
	}

	return &metricParams, nil
}
