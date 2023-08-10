package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadStats(t *testing.T) {
	defaultStats := Stats{}
	stats := Stats{}

	ReadStats(&stats)

	assert.NotEqual(t, defaultStats, stats)
}

func TestStats_GetGaugeMap(t *testing.T) {
	stats := Stats{}

	assert.Len(t, stats.GetGaugeMap(), 27)
}

func TestStats_GetCounterMap(t *testing.T) {
	stats := Stats{}

	assert.Len(t, stats.GetCounterMap(), 1)
}
