package metric

type Type string

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
)

func (t Type) ToString() string {
	return string(t)
}
