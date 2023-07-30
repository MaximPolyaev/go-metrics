package metric

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Name string

type Gauge struct {
	Name  Name
	Value float64
}

type Counter struct {
	Name  Name
	Value int
}
