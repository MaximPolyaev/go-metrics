package filestorage

import (
	"encoding/json"
	"errors"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"os"
	"sync"
)

type Storage struct {
	mu       sync.RWMutex
	filePath string
}

func New(filePath string) *Storage {
	return &Storage{filePath: filePath}
}

func (s *Storage) GetAll() ([]metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var mSlice []metric.Metric

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return mSlice, err
	}

	if len(data) == 0 {
		return mSlice, nil
	}

	if err := json.Unmarshal(data, &mSlice); err != nil {
		return nil, err
	}

	return mSlice, nil
}

func (s *Storage) SetAll(mSlice []metric.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(mSlice, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0666)
}
