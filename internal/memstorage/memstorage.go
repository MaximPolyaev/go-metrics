package memstorage

import "github.com/MaximPolyaev/go-metrics/internal/metric"

type MemStorage struct {
	values map[metric.Type]map[string]interface{}
}

func New() MemStorage {
	return MemStorage{values: make(map[metric.Type]map[string]interface{})}
}

func (s MemStorage) Set(mType metric.Type, key string, val interface{}) {
	if _, ok := s.values[mType]; !ok {
		s.values[mType] = make(map[string]interface{})
	}

	s.values[mType][key] = val
}

func (s MemStorage) Get(mType metric.Type, key string) (val interface{}, ok bool) {
	val, ok = s.values[mType][key]
	return
}

func (s MemStorage) GetAllByType(mType metric.Type) (values map[string]interface{}, ok bool) {
	values, ok = s.values[mType]

	return
}
