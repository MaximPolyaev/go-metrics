package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat64FromBytes(t *testing.T) {
	want := 2.2

	f := Float64FromBytes([]byte{154, 153, 153, 153, 153, 153, 1, 64})

	assert.Equal(t, want, f)
}

func TestFloat64ToByte(t *testing.T) {
	want := []byte{154, 153, 153, 153, 153, 153, 1, 64}

	b, err := Float64ToByte(2.2)

	assert.NoError(t, err)
	assert.Equal(t, want, b)
}

func TestIntToByte(t *testing.T) {
	i := IntToByte(10)

	assert.Equal(t, []byte{49, 48}, i)
}

func TestIntFromBytes(t *testing.T) {
	b, err := IntFromBytes([]byte{49, 48})

	assert.NoError(t, err)
	assert.Equal(t, 10, b)
}
