package defaultstats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats_ReadStats(t *testing.T) {
	defaultStats := New()
	stats := New()

	stats.ReadStats()

	assert.NotEqual(t, defaultStats, stats)
}

func TestStats_AsMetrics(t *testing.T) {
	stats := New()

	assert.Len(t, stats.AsMetrics(), 29)
}
