package memstorage

type MemStorage interface {
	Set(namespace string, key string, val interface{})
	Get(category string, key string) (val interface{}, ok bool)
	GetValuesByNamespace(namespace string) (values map[string]interface{}, ok bool)
}
type memStorage struct {
	values map[string]map[string]interface{}
}

func NewMemStorage() MemStorage {
	return memStorage{values: make(map[string]map[string]interface{})}
}

func (s memStorage) Set(namespace string, key string, val interface{}) {
	if _, ok := s.values[namespace]; !ok {
		s.values[namespace] = make(map[string]interface{})
	}

	s.values[namespace][key] = val
}

func (s memStorage) Get(namespace string, key string) (val interface{}, ok bool) {
	val, ok = s.values[namespace][key]
	return
}

func (s memStorage) GetValuesByNamespace(namespace string) (values map[string]interface{}, ok bool) {
	values, ok = s.values[namespace]

	return
}
