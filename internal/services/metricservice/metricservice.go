package metricservice

import (
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type MetricService struct {
	memStorage  memStorage
	fileStorage fileStorage
	storeCfg    *config.StoreConfig
	log         *logger.Logger
}

type memStorage interface {
	Set(mType metric.Type, val metric.Metric)
	Get(mType metric.Type, id string) (val metric.Metric, ok bool)
	GetAllByType(mType metric.Type) (values map[string]metric.Metric, ok bool)
}

type fileStorage interface {
	SetAll([]metric.Metric) error
	GetAll() ([]metric.Metric, error)
}

func New(
	memStorage memStorage,
	fileStorage fileStorage,
	storeCfg *config.StoreConfig,
	log *logger.Logger,
) (*MetricService, error) {
	ms := &MetricService{
		memStorage:  memStorage,
		fileStorage: fileStorage,
		storeCfg:    storeCfg,
		log:         log,
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

func (s *MetricService) Update(mm *metric.Metric) *metric.Metric {
	switch mm.MType {
	case metric.GaugeType:
		s.memStorage.Set(mm.MType, *mm)
	case metric.CounterType:
		existDelta, ok := s.memStorage.Get(mm.MType, mm.ID)

		if ok {
			*mm.Delta += *existDelta.Delta
		}

		s.memStorage.Set(mm.MType, *mm)
	}

	if s.storeCfg != nil && *s.storeCfg.StoreInterval == 0 {
		s.Sync()
	}

	return mm
}

func (s *MetricService) Get(mm *metric.Metric) (*metric.Metric, bool) {
	existMm, ok := s.memStorage.Get(mm.MType, mm.ID)

	if !ok {
		return mm, false
	}

	return &existMm, true
}

func (s *MetricService) GetAll() []metric.Metric {
	var mSlice []metric.Metric

	for _, mType := range metric.Types() {
		metricMap, ok := s.memStorage.GetAllByType(mType)
		if ok {
			for _, m := range metricMap {
				mSlice = append(mSlice, m)
			}
		}
	}

	return mSlice
}

func (s *MetricService) Sync() {
	if err := s.store(); err != nil {
		s.log.Error(err)
	}
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

func (s *MetricService) restore() error {
	if !*s.storeCfg.Restore {
		return nil
	}

	mSlice, err := s.fileStorage.GetAll()
	if err != nil {
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
	mSlice := s.GetAll()
	if mSlice == nil {
		return nil
	}

	return s.fileStorage.SetAll(mSlice)
}
