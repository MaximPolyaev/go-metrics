package memstorage

import "github.com/MaximPolyaev/go-metrics/internal/metric"

type MemStorage struct {
	values map[metric.Type]map[string]metric.Metric
}

func New() *MemStorage {
	return &MemStorage{values: make(map[metric.Type]map[string]metric.Metric)}
}

func (s *MemStorage) Set(mType metric.Type, val metric.Metric) {
	if _, ok := s.values[mType]; !ok {
		s.values[mType] = make(map[string]metric.Metric)
	}

	s.values[mType][val.ID] = val
}

func (s *MemStorage) Get(mType metric.Type, id string) (val metric.Metric, ok bool) {
	val, ok = s.values[mType][id]
	return
}

func (s *MemStorage) GetAllByType(mType metric.Type) (values map[string]metric.Metric, ok bool) {
	values, ok = s.values[mType]

	return
}
