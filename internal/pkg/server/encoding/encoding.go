package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

func Float64ToByte(f float64) ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.LittleEndian, f); err != nil {
		return buf.Bytes(), fmt.Errorf("binary.Write failed: %v", err)
	}

	return buf.Bytes(), nil
}

func IntToByte(i int) []byte {
	return []byte(strconv.Itoa(i))
}

func IntFromBytes(bytes []byte) (int, error) {
	s := string(bytes)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("can't convert this []byte to int: %v", err)
	}

	return i, nil
}

func Float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
