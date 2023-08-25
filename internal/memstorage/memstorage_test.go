package memstorage

import (
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

func Test_memStorage_Get(t *testing.T) {
	s := New()
	s.Set(metric.CounterType, "test key", 1)

	type args struct {
		mType metric.Type
		key   string
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
				key:   "test key",
			},
			wantVal: nil,
			wantOk:  false,
		},
		{
			name: "test case #2",
			args: args{
				mType: metric.CounterType,
				key:   "test key",
			},
			wantVal: 1,
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := s.Get(tt.args.mType, tt.args.key)
			assert.Equal(t, tt.wantVal, gotVal)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}

func Test_memStorage_GetValuesByNamespace(t *testing.T) {
	s := New()
	s.Set(metric.CounterType, "test key", 1)

	tests := []struct {
		name       string
		mType      metric.Type
		wantValues map[string]interface{}
		wantOk     bool
	}{
		{
			name:       "test case #1",
			mType:      metric.GaugeType,
			wantValues: map[string]interface{}(nil),
			wantOk:     false,
		},
		{
			name:  "test case #2",
			mType: metric.CounterType,
			wantValues: map[string]interface{}{
				"test key": 1,
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues, gotOk := s.GetAllByType(tt.mType)
			assert.Equal(t, tt.wantValues, gotValues)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
