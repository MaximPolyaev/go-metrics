package services

import (
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/stretchr/testify/assert"
)

func Test_counterService_Update(t *testing.T) {
	storage := memstorage.NewMemStorage()

	type args struct {
		name   string
		valStr string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test case #1",
			args: args{
				name:   "some",
				valStr: "1",
			},
			wantErr: false,
		},
		{
			name: "test case #2",
			args: args{
				name:   "some",
				valStr: "0",
			},
			wantErr: false,
		},
		{
			name: "test case #3",
			args: args{
				name:   "some",
				valStr: "-1",
			},
			wantErr: false,
		},
		{
			name: "test case #4",
			args: args{
				name:   "some",
				valStr: "1.1",
			},
			wantErr: true,
		},
		{
			name: "test case #5",
			args: args{
				name:   "some",
				valStr: "",
			},
			wantErr: true,
		},
		{
			name: "test case #6",
			args: args{
				name:   "",
				valStr: "1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateService := &counterService{
				s: storage,
			}

			err := updateService.Update(tt.args.name, tt.args.valStr)

			if tt.wantErr {
				assert.Error(t, err)
			}

			if err := updateService.Update(tt.args.name, tt.args.valStr); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_counterService_GetValues(t *testing.T) {
	tests := []struct {
		name    string
		storage memstorage.MemStorage
		want    map[string]string
	}{
		{
			name:    "empty storage",
			storage: memstorage.NewMemStorage(),
			want:    map[string]string{},
		},
		{
			name:    "not empty storage",
			storage: MakeStorageWithCounterValue(),
			want: map[string]string{
				"test": "10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateService := &counterService{
				s: tt.storage,
			}

			got, err := updateService.GetValues()

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_counterService_GetValue(t *testing.T) {
	tests := []struct {
		name    string
		storage memstorage.MemStorage
		mName   string
		ok      bool
		want    string
		wantErr bool
	}{
		{
			name:    "empty storage",
			storage: memstorage.NewMemStorage(),
			mName:   "not exist",
			ok:      false,
			want:    "",
			wantErr: true,
		},
		{
			name:    "not empty storage",
			storage: MakeStorageWithCounterValue(),
			mName:   "test",
			ok:      true,
			want:    "10",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mService := &counterService{
				s: tt.storage,
			}

			got, ok, err := mService.GetValue(tt.mName)

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

func MakeStorageWithCounterValue() memstorage.MemStorage {
	storage := memstorage.NewMemStorage()
	storage.Set(string(metric.CounterType), "test", encoding.IntToByte(10))

	return storage
}
