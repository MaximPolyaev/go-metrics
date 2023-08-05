package metric

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
)
