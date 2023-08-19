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

func TestStats_AsMetrics(t *testing.T) {
	stats := Stats{}

	assert.Len(t, stats.AsMetrics(), 29)
}
