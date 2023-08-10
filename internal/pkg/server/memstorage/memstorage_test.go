package memstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_memStorage_Get(t *testing.T) {
	s := New()
	s.Set("ns test", "test key", 1)

	type args struct {
		namespace string
		key       string
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
				namespace: "ns test",
				key:       "test key",
			},
			wantVal: 1,
			wantOk:  true,
		},
		{
			name: "test case #2",
			args: args{
				namespace: "ns test 2",
				key:       "test key",
			},
			wantVal: nil,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := s.Get(tt.args.namespace, tt.args.key)
			assert.Equal(t, tt.wantVal, gotVal)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}

func Test_memStorage_GetValuesByNamespace(t *testing.T) {
	s := New()
	s.Set("ns test", "test key", 1)

	tests := []struct {
		name       string
		namespace  string
		wantValues map[string]interface{}
		wantOk     bool
	}{
		{
			name:       "test case #1",
			namespace:  "not exists",
			wantValues: map[string]interface{}(nil),
			wantOk:     false,
		},
		{
			name:      "test case #2",
			namespace: "ns test",
			wantValues: map[string]interface{}{
				"test key": 1,
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues, gotOk := s.GetValuesByNamespace(tt.namespace)
			assert.Equal(t, tt.wantValues, gotValues)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}
