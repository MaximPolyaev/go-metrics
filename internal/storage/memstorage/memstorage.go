package memstorage

import (
	"context"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"sync"
)

type Storage struct {
	mu sync.RWMutex

	values map[metric.Type]map[string]metric.Metric
}

func New() *Storage {
	return &Storage{values: make(map[metric.Type]map[string]metric.Metric)}
}

func (s *Storage) Set(_ context.Context, mType metric.Type, val metric.Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.values[mType]; !ok {
		s.values[mType] = make(map[string]metric.Metric)
	}

	s.values[mType][val.ID] = val
}

func (s *Storage) BatchSet(ctx context.Context, mSlice []metric.Metric) {
	for _, m := range mSlice {
		s.Set(ctx, m.MType, m)
	}
}

func (s *Storage) Get(_ context.Context, mType metric.Type, id string) (val metric.Metric, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok = s.values[mType][id]
	return
}

func (s *Storage) GetAllByType(_ context.Context, mType metric.Type) (values map[string]metric.Metric, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	values, ok = s.values[mType]

	return
}
