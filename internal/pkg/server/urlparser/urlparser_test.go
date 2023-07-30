package urlparser

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeGaugeMetricByURLPath(t *testing.T) {
	tests := []struct {
		name    string
		urlPath string
		want    *metric.Gauge
		wantErr bool
	}{
		{
			name:    "negative empty url case #1",
			urlPath: "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #2",
			urlPath: " ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #3",
			urlPath: "/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #4",
			urlPath: "//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #1",
			urlPath: "test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #2",
			urlPath: "/test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #3",
			urlPath: "/test/    /",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #4",
			urlPath: "/    /test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case some text value",
			urlPath: "test/test",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case float value",
			urlPath: "test/3,3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "positive url case #1",
			urlPath: "test/0",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #2",
			urlPath: "1234/0",
			want:    &metric.Gauge{Name: metric.Name("1234"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #3",
			urlPath: "_/0",
			want:    &metric.Gauge{Name: metric.Name("_"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #4",
			urlPath: "test/-0",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #5",
			urlPath: "test/-1",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: -1},
			wantErr: false,
		},
		{
			name:    "positive url case #6",
			urlPath: "test/-1.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: -1.5},
			wantErr: false,
		},
		{
			name:    "positive url case #7",
			urlPath: "test/-01.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: -1.5},
			wantErr: false,
		},
		{
			name:    "positive url case #8",
			urlPath: "test/-00.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: -0.5},
			wantErr: false,
		},
		{
			name:    "positive url case #9",
			urlPath: "test/1",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 1},
			wantErr: false,
		},
		{
			name:    "positive url case #10",
			urlPath: "test/1.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 1.5},
			wantErr: false,
		},
		{
			name:    "positive url case #11",
			urlPath: "test/01.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 1.5},
			wantErr: false,
		},
		{
			name:    "positive url case #12",
			urlPath: "test/00.5",
			want:    &metric.Gauge{Name: metric.Name("test"), Value: 0.5},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeGaugeMetricByURLPath(tt.urlPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMakeCounterMetricByURLPath(t *testing.T) {
	tests := []struct {
		name    string
		urlPath string
		want    *metric.Counter
		wantErr bool
	}{
		{
			name:    "negative empty url case #1",
			urlPath: "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #2",
			urlPath: " ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #3",
			urlPath: "/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative empty url case #4",
			urlPath: "//",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #1",
			urlPath: "test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #2",
			urlPath: "/test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #3",
			urlPath: "/test/    /",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case #4",
			urlPath: "/    /test/",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case some text value",
			urlPath: "test/test",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case float value #1",
			urlPath: "test/3,3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case float value #2",
			urlPath: "test/3.3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative url case float value #3",
			urlPath: "test/-3.3",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "positive url case #1",
			urlPath: "test/0",
			want:    &metric.Counter{Name: metric.Name("test"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #2",
			urlPath: "1234/0",
			want:    &metric.Counter{Name: metric.Name("1234"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #3",
			urlPath: "_/0",
			want:    &metric.Counter{Name: metric.Name("_"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #4",
			urlPath: "test/-0",
			want:    &metric.Counter{Name: metric.Name("test"), Value: 0},
			wantErr: false,
		},
		{
			name:    "positive url case #5",
			urlPath: "test/-1",
			want:    &metric.Counter{Name: metric.Name("test"), Value: -1},
			wantErr: false,
		},
		{
			name:    "positive url case #6",
			urlPath: "test/1",
			want:    &metric.Counter{Name: metric.Name("test"), Value: 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeCounterMetricByURLPath(tt.urlPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
