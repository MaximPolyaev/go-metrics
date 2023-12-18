// Package hash используется для кодирования данных в sha256
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encode - encode []byte value to hash string key
func Encode(value []byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	if _, err := h.Write(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
