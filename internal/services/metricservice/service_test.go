package metricservice

import (
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

type mockMemStorage struct{}

func TestMetricService_Update(t *testing.T) {
	type args struct {
		name       string
		metricType metric.Type
		valStr     string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "counter test case #1",
			args: args{
				name:       "some",
				valStr:     "1",
				metricType: metric.CounterType,
			},
			wantErr: false,
		},
		{
			name: "counter test case #2",
			args: args{
				name:       "some",
				valStr:     "0",
				metricType: metric.CounterType,
			},
			wantErr: false,
		},
		{
			name: "counter test case #3",
			args: args{
				name:       "some",
				valStr:     "-1",
				metricType: metric.CounterType,
			},
			wantErr: false,
		},
		{
			name: "counter test case #4",
			args: args{
				name:       "some",
				valStr:     "1.1",
				metricType: metric.CounterType,
			},
			wantErr: true,
		},
		{
			name: "counter test case #5",
			args: args{
				name:       "some",
				valStr:     "",
				metricType: metric.CounterType,
			},
			wantErr: true,
		},
		{
			name: "counter test case #6",
			args: args{
				name:       "",
				valStr:     "1",
				metricType: metric.CounterType,
			},
			wantErr: true,
		},
		{
			name: "gauge test case #1",
			args: args{
				name:       "some",
				valStr:     "1",
				metricType: metric.GaugeType,
			},
			wantErr: false,
		},
		{
			name: "gauge test case #2",
			args: args{
				name:       "some",
				valStr:     "0",
				metricType: metric.GaugeType,
			},
			wantErr: false,
		},
		{
			name: "gauge test case #3",
			args: args{
				name:       "some",
				valStr:     "-1",
				metricType: metric.GaugeType,
			},
			wantErr: false,
		},
		{
			name: "gauge test case #4",
			args: args{
				name:       "some",
				valStr:     "1.1",
				metricType: metric.GaugeType,
			},
			wantErr: false,
		},
		{
			name: "gauge test case #5",
			args: args{
				name:       "some",
				valStr:     "",
				metricType: metric.GaugeType,
			},
			wantErr: true,
		},
		{
			name: "gauge test case #6",
			args: args{
				name:       "",
				valStr:     "1",
				metricType: metric.GaugeType,
			},
			wantErr: true,
		},
		{
			name: "gauge test case #7",
			args: args{
				name:       "some",
				valStr:     "2,2",
				metricType: metric.GaugeType,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(mockMemStorage{})

			err := s.Update(tt.args.metricType, tt.args.name, tt.args.valStr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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

	s := New(mockMemStorage{})

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
		metricName string
		metricType metric.Type
		ok         bool
		want       string
		wantErr    bool
	}{
		{
			name:       "empty counter value",
			metricName: "not exist",
			metricType: metric.CounterType,
			ok:         false,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "not empty counter value",
			metricName: "test",
			metricType: metric.CounterType,
			ok:         true,
			want:       "10",
			wantErr:    false,
		},
		{
			name:       "empty gauge value",
			metricName: "not exist",
			metricType: metric.GaugeType,
			ok:         false,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "not empty gauge value",
			metricName: "test",
			metricType: metric.GaugeType,
			ok:         true,
			want:       "1.1",
			wantErr:    false,
		},
	}

	s := New(mockMemStorage{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := s.GetValue(tt.metricType, tt.metricName)

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

func (s mockMemStorage) Set(_ string, _ string, _ interface{}) {}
func (s mockMemStorage) Get(namespace string, key string) (val interface{}, ok bool) {
	switch namespace {
	case metric.CounterType.ToString():
		if key == "test" {
			return 10, true
		}
	case metric.GaugeType.ToString():
		if key == "test" {
			return 1.1, true
		}
	}
	return nil, false
}

func (s mockMemStorage) GetValuesByNamespace(namespace string) (values map[string]interface{}, ok bool) {
	switch namespace {
	case metric.CounterType.ToString():
		return map[string]interface{}{
			"test": 10,
		}, true
	case metric.GaugeType.ToString():
		return map[string]interface{}{
			"test": 1.1,
		}, true
	}

	return nil, false
}
