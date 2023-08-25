package metricservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type MetricService struct {
	storage  memStorage
	storeCfg *config.StoreConfig
	log      *logger.Logger
}

type memStorage interface {
	Set(mType metric.Type, key string, val interface{})
	Get(mType metric.Type, key string) (val interface{}, ok bool)
	GetAllByType(mType metric.Type) (values map[string]interface{}, ok bool)
}

func New(s memStorage, storeCfg *config.StoreConfig, log *logger.Logger) (*MetricService, error) {
	ms := &MetricService{
		storage:  s,
		storeCfg: storeCfg,
		log:      log,
	}

	if storeCfg != nil {
		if err := ms.restore(); err != nil {
			return nil, err
		}

		if *storeCfg.StoreInterval != 0 {
			ms.async()
		}
	}

	return ms, nil
}

func (s *MetricService) Update(mm *metric.Metrics) *metric.Metrics {
	switch mm.MType {
	case metric.GaugeType:
		s.storage.Set(mm.MType, mm.ID, *mm.Value)
	case metric.CounterType:
		existValue, ok := s.storage.Get(mm.MType, mm.ID)

		if ok {
			*mm.Delta += existValue.(int64)
		}

		s.storage.Set(mm.MType, mm.ID, *mm.Delta)
	}

	s.sync()

	return mm
}

func (s *MetricService) Get(mm *metric.Metrics) *metric.Metrics {
	mm.ValueInit()

	value, ok := s.storage.Get(mm.MType, mm.ID)

	if !ok {
		return mm
	}

	switch mm.MType {
	case metric.GaugeType:
		*mm.Value = value.(float64)
	case metric.CounterType:
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

	return map[string]string(nil), errors.New("unexpected metricservice type: " + mType.ToString())
}

func (s *MetricService) GetValue(mType metric.Type, name string) (value string, ok bool, err error) {
	switch mType {
	case metric.GaugeType:
		return s.getGaugeValue(name)
	case metric.CounterType:
		return s.getCounterValue(name)
	}

	return "", false, errors.New("unexpected metricservice type: " + mType.ToString())
}

func (s *MetricService) getGaugeValues() (map[string]string, error) {
	values, ok := s.storage.GetAllByType(metric.GaugeType)

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
	values, ok := s.storage.GetAllByType(metric.CounterType)

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
	value, ok := s.storage.Get(metric.GaugeType, name)

	if !ok {
		return "", ok, errors.New("metricservice " + name + " not found")
	}

	strValue = fmt.Sprintf("%g", value.(float64))

	return
}

func (s *MetricService) getCounterValue(name string) (strValue string, ok bool, err error) {
	value, ok := s.storage.Get(metric.CounterType, name)

	if !ok {
		return "", ok, errors.New("metricservice " + name + " not found")
	}

	strValue = strconv.Itoa(int(value.(int64)))

	return
}

func (s *MetricService) async() {
	storeInterval := time.NewTicker(time.Duration(*s.storeCfg.StoreInterval) * time.Second)

	go func() {
		for {
			<-storeInterval.C

			if err := s.store(); err != nil {
				s.log.Error(err)
			}
		}
	}()
}

func (s *MetricService) sync() {
	if s.storeCfg == nil || *s.storeCfg.StoreInterval != 0 {
		return
	}

	if err := s.store(); err != nil {
		s.log.Error(err)
	}
}

func (s *MetricService) restore() error {
	if !*s.storeCfg.Restore {
		return nil
	}

	data, err := os.ReadFile(*s.storeCfg.FileStoragePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	var mSlice []metric.Metrics

	if err := json.Unmarshal(data, &mSlice); err != nil {
		return err
	}

	if mSlice == nil {
		return nil
	}

	for _, m := range mSlice {
		s.Update(&m)
	}

	return nil
}

func (s *MetricService) store() error {
	mSlice := s.getAll()
	if mSlice == nil {
		return nil
	}

	data, err := json.MarshalIndent(mSlice, "", " ")
	if err != nil {
		return nil
	}

	return os.WriteFile(*s.storeCfg.FileStoragePath, data, 0666)
}

func (s *MetricService) getAll() []metric.Metrics {
	values, ok := s.storage.GetAllByType(metric.CounterType)

	var mSlice []metric.Metrics

	if ok {
		for k, v := range values {
			tmpV := v.(int64)

			mSlice = append(mSlice, metric.Metrics{
				ID:    k,
				MType: metric.CounterType,
				Delta: &tmpV,
			})
		}
	}

	values, ok = s.storage.GetAllByType(metric.GaugeType)

	if ok {
		for k, v := range values {
			tmpV := v.(float64)

			mSlice = append(mSlice, metric.Metrics{
				ID:    k,
				MType: metric.GaugeType,
				Value: &tmpV,
			})
		}
	}

	return mSlice
}
