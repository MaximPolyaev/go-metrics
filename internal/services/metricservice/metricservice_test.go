package metricservice

import (
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

type mockMemStorage struct{}

func TestMetricService_Update(t *testing.T) {
	type args struct {
		id         string
		metricType metric.Type
		delta      int64
		value      float64
	}

	tests := []struct {
		name      string
		args      args
		wantDelta int64
		wantValue float64
	}{
		{
			name: "counter test case #1",
			args: args{
				id:         "test",
				delta:      1,
				metricType: metric.CounterType,
			},
			wantDelta: 11,
		},
		{
			name: "counter test case #2",
			args: args{
				id:         "test",
				delta:      0,
				metricType: metric.CounterType,
			},
			wantDelta: 10,
		},
		{
			name: "counter test case #3",
			args: args{
				id:         "test",
				delta:      -1,
				metricType: metric.CounterType,
			},
			wantDelta: 9,
		},
		{
			name: "gauge test case #1",
			args: args{
				id:         "test",
				value:      1,
				metricType: metric.GaugeType,
			},
			wantValue: 1,
		},
		{
			name: "gauge test case #2",
			args: args{
				id:         "test",
				value:      0,
				metricType: metric.GaugeType,
			},
			wantValue: 0,
		},
		{
			name: "gauge test case #3",
			args: args{
				id:         "test",
				value:      -1,
				metricType: metric.GaugeType,
			},
			wantValue: -1,
		},
		{
			name: "gauge test case #4",
			args: args{
				id:         "test",
				value:      1.1,
				metricType: metric.GaugeType,
			},
			wantValue: 1.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(mockMemStorage{}, nil, nil)
			assert.NoError(t, err)

			mm := metric.Metrics{
				ID:    tt.args.id,
				MType: tt.args.metricType,
			}

			switch mm.MType {
			case metric.GaugeType:
				mm.Value = &tt.args.value
			case metric.CounterType:
				mm.Delta = &tt.args.delta
			}

			mm = *s.Update(&mm)

			switch mm.MType {
			case metric.GaugeType:
				assert.NotNil(t, mm.Value)
				assert.Equal(t, tt.wantValue, *mm.Value)
			case metric.CounterType:
				assert.NotNil(t, mm.Delta)
				assert.Equal(t, tt.wantDelta, *mm.Delta)
			}
		})
	}
}

func TestMetricService_GetValues(t *testing.T) {
	tests := []struct {
		name       string
		metricType metric.Type
		want       map[string]string
	}{
		{
			name:       "counter values",
			metricType: metric.CounterType,
			want: map[string]string{
				"test": "10",
			},
		},
		{
			name:       "gauge values",
			metricType: metric.GaugeType,
			want: map[string]string{
				"test": "1.1",
			},
		},
	}

	s, err := New(mockMemStorage{}, nil, nil)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetValues(tt.metricType)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetricService_GetValue(t *testing.T) {
	tests := []struct {
		name       string
		metricId   string
		metricType metric.Type
		ok         bool
		want       string
		wantErr    bool
	}{
		{
			name:       "empty counter value",
			metricId:   "not exist",
			metricType: metric.CounterType,
			ok:         false,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "not empty counter value",
			metricId:   "test",
			metricType: metric.CounterType,
			ok:         true,
			want:       "10",
			wantErr:    false,
		},
		{
			name:       "empty gauge value",
			metricId:   "not exist",
			metricType: metric.GaugeType,
			ok:         false,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "not empty gauge value",
			metricId:   "test",
			metricType: metric.GaugeType,
			ok:         true,
			want:       "1.1",
			wantErr:    false,
		},
	}

	s, err := New(mockMemStorage{}, nil, nil)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := s.GetValue(tt.metricType, tt.metricId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func (s mockMemStorage) Set(_ metric.Type, _ string, _ interface{}) {}
func (s mockMemStorage) Get(mType metric.Type, id string) (val interface{}, ok bool) {
	switch mType {
	case metric.CounterType:
		if id == "test" {
			return int64(10), true
		}
	case metric.GaugeType:
		if id == "test" {
			return 1.1, true
		}
	}
	return nil, false
}

func (s mockMemStorage) GetAllByType(mType metric.Type) (values map[string]interface{}, ok bool) {
	switch mType {
	case metric.CounterType:
		return map[string]interface{}{
			"test": int64(10),
		}, true
	case metric.GaugeType:
		return map[string]interface{}{
			"test": 1.1,
		}, true
	}

	return nil, false
}
