package memstorage

type MemStorage struct {
	values map[string]map[string]interface{}
}

func New() MemStorage {
	return MemStorage{values: make(map[string]map[string]interface{})}
}

func (s MemStorage) Set(namespace string, key string, val interface{}) {
	if _, ok := s.values[namespace]; !ok {
		s.values[namespace] = make(map[string]interface{})
	}

	s.values[namespace][key] = val
}

func (s MemStorage) Get(namespace string, key string) (val interface{}, ok bool) {
	val, ok = s.values[namespace][key]
	return
}

func (s MemStorage) GetValuesByNamespace(namespace string) (values map[string]interface{}, ok bool) {
	values, ok = s.values[namespace]

	return
}
