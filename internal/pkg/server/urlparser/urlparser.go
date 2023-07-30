package urlparser

import (
	"errors"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"strconv"
	"strings"
)

func MakeGaugeMetricByURLPath(urlPath string) (*metric.Gauge, error) {
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

	return &metric.Gauge{
		Name:  metric.Name(name),
		Value: value,
	}, nil
}

func MakeCounterMetricByURLPath(urlPath string) (*metric.Counter, error) {
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

	return &metric.Counter{
		Name:  metric.Name(name),
		Value: value,
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
