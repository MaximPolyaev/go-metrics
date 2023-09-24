package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	assert.Equal(t,
		"02afb56304902c656fcb737cdd03de6205bb6d401da2812efd9b2d36a08af159",
		Encode([]byte("test"), "key"),
	)
}
