package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

type Decoder struct {
	pk   *rsa.PrivateKey
	size int
}

func NewCryptoDecoder(pk *rsa.PrivateKey) *Decoder {
	return &Decoder{pk: pk, size: pk.Size()}
}

func (d *Decoder) Decode(data []byte) ([]byte, error) {
	encodedChunks := dataToChunks(data, d.size)
	var buf bytes.Buffer

	for _, encodedChunk := range encodedChunks {
		decryptChunk, err := rsa.DecryptPKCS1v15(rand.Reader, d.pk, encodedChunk)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(decryptChunk)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
