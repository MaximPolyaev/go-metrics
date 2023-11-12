package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encode - encode []byte value to hash string key
func Encode(value []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(value)

	return hex.EncodeToString(h.Sum(nil))
}
