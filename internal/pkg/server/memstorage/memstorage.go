package memstorage

type MemStorage interface {
	Set(namespace string, key string, val []byte)
	Get(category string, key string) (val []byte, ok bool)
	GetValuesByNamespace(namespace string) (values map[string][]byte, ok bool)
}

type memStorage struct {
	values map[string]map[string][]byte
}

func NewMemStorage() MemStorage {
	return memStorage{values: make(map[string]map[string][]byte)}
}

func (s memStorage) Set(namespace string, key string, val []byte) {
	if _, ok := s.values[namespace]; !ok {
		s.values[namespace] = make(map[string][]byte)
	}

	s.values[namespace][key] = val
}

func (s memStorage) Get(namespace string, key string) (val []byte, ok bool) {
	val, ok = s.values[namespace][key]
	return
}

func (s memStorage) GetValuesByNamespace(namespace string) (values map[string][]byte, ok bool) {
	values, ok = s.values[namespace]

	return
}
