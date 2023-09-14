package memstorage

import (
	"context"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

func Test_memStorage_Get(t *testing.T) {
	s := New()

	delta := int64(1)

	wantVal := metric.Metric{
		ID:    "test id",
		MType: metric.CounterType,
		Delta: &delta,
	}

	s.Set(context.TODO(), metric.CounterType, wantVal)

	type args struct {
		mType metric.Type
		id    string
	}

	tests := []struct {
		name    string
		args    args
		wantVal interface{}
		wantOk  bool
	}{
		{
			name: "test case #1",
			args: args{
				mType: metric.GaugeType,
				id:    "test id",
			},
			wantVal: metric.Metric{},
			wantOk:  false,
		},
		{
			name: "test case #2",
			args: args{
				mType: metric.CounterType,
				id:    "test id",
			},
			wantVal: wantVal,
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := s.Get(context.TODO(), tt.args.mType, tt.args.id)
			assert.Equal(t, tt.wantVal, gotVal)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}

func Test_memStorage_GetValuesByNamespace(t *testing.T) {
	s := New()

	delta := int64(1)

	wantVal := metric.Metric{
		ID:    "test id",
		MType: metric.CounterType,
		Delta: &delta,
	}

	s.Set(context.TODO(), metric.CounterType, wantVal)

	tests := []struct {
		name       string
		mType      metric.Type
		wantValues map[string]metric.Metric
		wantOk     bool
	}{
		{
			name:       "test case #1",
			mType:      metric.GaugeType,
			wantValues: map[string]metric.Metric(nil),
			wantOk:     false,
		},
		{
			name:  "test case #2",
			mType: metric.CounterType,
			wantValues: map[string]metric.Metric{
				"test id": wantVal,
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues, gotOk := s.GetAllByType(context.TODO(), tt.mType)
			assert.Equal(t, tt.wantValues, gotValues)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
