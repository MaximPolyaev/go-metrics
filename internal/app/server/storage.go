package server

type MemStorageInterface interface {
	writeGaugeMetric(metric gaugeMetric)
	writeCounterMetric(metric counterMetric)
}

type MemStorage struct {
	gaugeMetrics   map[string]gaugeMetric
	counterMetrics map[string]counterMetric
}

var memStorage MemStorage

func init() {
	memStorage = MemStorage{
		gaugeMetrics:   make(map[string]gaugeMetric),
		counterMetrics: make(map[string]counterMetric),
	}
}

func getMemStorage() MemStorage {
	return memStorage
}

func (m MemStorage) writeGaugeMetric(metric gaugeMetric) {
	m.gaugeMetrics[metric.name] = metric
}

func (m MemStorage) writeCounterMetric(metric counterMetric) {
	existMetric, ok := m.counterMetrics[metric.name]

	if ok {
		existMetric.value += metric.value
		m.counterMetrics[metric.name] = existMetric
		return
	}

	m.counterMetrics[metric.name] = metric
}
