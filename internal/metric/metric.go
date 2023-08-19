package metric

import (
	"errors"
	"fmt"
)

type Type string

type Metrics struct {
	ID    string   `json:"id"`
	MType Type     `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

const (
	GaugeType   = Type("gauge")
	CounterType = Type("counter")
)

func (t Type) ToString() string {
	return string(t)
}

func (t Type) Validate() error {
	switch t {
	case GaugeType:
		return nil
	case CounterType:
		return nil
	}

	return errors.New("invalid metric type: " + t.ToString())
}

func (m *Metrics) ValueInit() {
	switch m.MType {
	case CounterType:
		m.Delta = new(int64)
	case GaugeType:
		m.Value = new(float64)
	}
}

func (m *Metrics) Validate() error {
	if len(m.ID) == 0 {
		return errors.New("metric ID must be not empty")
	}

	if err := m.MType.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *Metrics) ValidateWithValue() error {
	if err := m.Validate(); err != nil {
		return err
	}

	switch m.MType {
	case CounterType:
		if m.Delta == nil {
			return fmt.Errorf("empty value for metric %s type", m.MType.ToString())
		}
	case GaugeType:
		if m.Value == nil {
			return fmt.Errorf("empty value for metric %s type", m.MType.ToString())
		}
	}

	return nil
}
