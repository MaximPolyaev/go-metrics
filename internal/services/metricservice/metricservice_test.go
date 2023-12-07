package metricservice

import (
	"context"
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
			s, err := New(mockMemStorage{}, nil, nil, nil)
			assert.NoError(t, err)

			mm := metric.Metric{
				ID:    tt.args.id,
				MType: tt.args.metricType,
			}

			switch mm.MType {
			case metric.GaugeType:
				mm.Value = &tt.args.value
			case metric.CounterType:
				mm.Delta = &tt.args.delta
			}

			mm = *s.Update(context.TODO(), &mm)

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

func BenchmarkMetricService_Update(b *testing.B) {
	b.Run("gauge", func(b *testing.B) {
		s, err := New(mockMemStorage{}, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		mm := metric.Metric{
			ID:    "test id",
			MType: metric.GaugeType,
			Value: new(float64),
		}
		*mm.Value = 2

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Update(context.TODO(), &mm)
		}
	})

	b.Run("counter", func(b *testing.B) {
		s, err := New(mockMemStorage{}, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		mm := metric.Metric{
			ID:    "test id",
			MType: metric.CounterType,
			Delta: new(int64),
		}
		*mm.Delta = 2
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Update(context.TODO(), &mm)
		}
	})
}

func (s mockMemStorage) Set(_ context.Context, _ metric.Type, _ metric.Metric) {}
func (s mockMemStorage) Get(ctx context.Context, mType metric.Type, id string) (val metric.Metric, ok bool) {
	metricMap, ok := s.GetAllByType(ctx, mType)
	if !ok {
		return
	}

	val, ok = metricMap[id]

	return
}

func (s mockMemStorage) GetAllByType(_ context.Context, mType metric.Type) (values map[string]metric.Metric, ok bool) {
	var delta int64
	var value float64

	delta = 10
	value = 1.1

	switch mType {
	case metric.CounterType:
		return map[string]metric.Metric{
			"test": {ID: "test", MType: metric.CounterType, Delta: &delta},
		}, true
	case metric.GaugeType:
		return map[string]metric.Metric{
			"test": {ID: "test", MType: metric.GaugeType, Value: &value},
		}, true
	}

	return nil, false
}

func (s mockMemStorage) BatchSet(_ context.Context, _ []metric.Metric) {}
