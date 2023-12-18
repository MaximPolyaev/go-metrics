// Package metricservice work with metric data
package metricservice

import (
	"context"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type MetricService struct {
	mStorage    metricStorage
	fileStorage fileStorage
	serverCfg   *config.ServerConfig
	log         *logger.Logger
}

type metricStorage interface {
	Set(ctx context.Context, mType metric.Type, val metric.Metric)
	Get(ctx context.Context, mType metric.Type, id string) (val metric.Metric, ok bool)
	BatchSet(ctx context.Context, mSlice []metric.Metric)
	GetAll(ctx context.Context) []metric.Metric
}

type fileStorage interface {
	SetAll([]metric.Metric) error
	GetAll() ([]metric.Metric, error)
}

func New(
	memStorage metricStorage,
	fileStorage fileStorage,
	serverCfg *config.ServerConfig,
	log *logger.Logger,
) (*MetricService, error) {
	ms := &MetricService{
		mStorage:    memStorage,
		fileStorage: fileStorage,
		serverCfg:   serverCfg,
		log:         log,
	}

	if serverCfg != nil {
		if err := ms.restore(); err != nil {
			return nil, err
		}

		if *serverCfg.StoreInterval != 0 {
			ms.async()
		}
	}

	return ms, nil
}

// Update - update metric with storage
func (s *MetricService) Update(ctx context.Context, mm *metric.Metric) *metric.Metric {
	switch mm.MType {
	case metric.GaugeType:
		s.mStorage.Set(context.Background(), mm.MType, *mm)
	case metric.CounterType:
		existDelta, ok := s.mStorage.Get(ctx, mm.MType, mm.ID)

		if ok {
			*mm.Delta += *existDelta.Delta
		}

		s.mStorage.Set(ctx, mm.MType, *mm)
	}

	if s.serverCfg != nil && *s.serverCfg.StoreInterval == 0 {
		s.Sync(context.Background())
	}

	return mm
}

// BatchUpdate - batch update metrics with recalculate metrics arg
func (s *MetricService) BatchUpdate(ctx context.Context, mSlice []metric.Metric) error {
	if len(mSlice) == 0 {
		return nil
	}

	tmpMetricMap := make(map[string]metric.Metric, len(mSlice))

	for _, m := range mSlice {
		if err := m.ValidateWithValue(); err != nil {
			return err
		}

		tmpKey := m.ID + "#" + m.MType.ToString()

		if m.MType == metric.CounterType {
			existCounter, ok := tmpMetricMap[tmpKey]

			if ok {
				*existCounter.Delta += *m.Delta
				continue
			}
		}

		tmpMetricMap[tmpKey] = m
	}

	updSlice := make([]metric.Metric, 0, len(tmpMetricMap))

	for _, m := range tmpMetricMap {
		if m.MType == metric.CounterType {
			existM, ok := s.mStorage.Get(ctx, m.MType, m.ID)

			if ok {
				*m.Delta += *existM.Delta
			}
		}

		updSlice = append(updSlice, m)
	}

	s.mStorage.BatchSet(ctx, updSlice)

	return nil
}

// Get - get metric from storage
func (s *MetricService) Get(ctx context.Context, mm *metric.Metric) (*metric.Metric, bool) {
	existMm, ok := s.mStorage.Get(ctx, mm.MType, mm.ID)

	if !ok {
		return mm, false
	}

	return &existMm, true
}

// GetAll - получить все метрики из хранилища
func (s *MetricService) GetAll(ctx context.Context) []metric.Metric {
	return s.mStorage.GetAll(ctx)
}

// Sync - sync all storage metrics with file storage
func (s *MetricService) Sync(ctx context.Context) {
	if err := s.store(ctx); err != nil {
		s.log.Error(err)
	}
}

func (s *MetricService) async() {
	storeInterval := time.NewTicker(time.Duration(*s.serverCfg.StoreInterval) * time.Second)

	go func() {
		for {
			<-storeInterval.C

			if err := s.store(context.Background()); err != nil {
				s.log.Error(err)
			}
		}
	}()
}

func (s *MetricService) restore() error {
	if !*s.serverCfg.Restore {
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
		s.Update(context.Background(), &m)
	}

	return nil
}

func (s *MetricService) store(ctx context.Context) error {
	mSlice := s.GetAll(ctx)
	if mSlice == nil {
		return nil
	}

	return s.fileStorage.SetAll(mSlice)
}
