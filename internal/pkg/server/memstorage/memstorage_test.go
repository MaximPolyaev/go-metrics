package memstorage

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_memStorage_Get(t *testing.T) {
	testBinaryInt := encoding.IntToByte(1)

	s := NewMemStorage()
	s.Set("ns test", "test key", testBinaryInt)

	type args struct {
		namespace string
		key       string
	}

	tests := []struct {
		name    string
		args    args
		wantVal []byte
		wantOk  bool
	}{
		{
			name: "test case #1",
			args: args{
				namespace: "ns test",
				key:       "test key",
			},
			wantVal: testBinaryInt,
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
	testBinaryInt := encoding.IntToByte(1)
	s := NewMemStorage()
	s.Set("ns test", "test key", testBinaryInt)

	tests := []struct {
		name       string
		namespace  string
		wantValues map[string][]byte
		wantOk     bool
	}{
		{
			name:       "test case #1",
			namespace:  "not exists",
			wantValues: map[string][]byte(nil),
			wantOk:     false,
		},
		{
			name:      "test case #2",
			namespace: "ns test",
			wantValues: map[string][]byte{
				"test key": testBinaryInt,
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
