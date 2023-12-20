package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

type Encoder struct {
	pk   *rsa.PublicKey
	size int
}

func NewCryptoEncoder(pk *rsa.PublicKey) *Encoder {
	return &Encoder{pk: pk, size: pk.Size() - 11}
}

func (e *Encoder) Encode(data []byte) ([]byte, error) {
	chunks := dataToChunks(data, e.size)
	buf := bytes.NewBuffer([]byte{})

	for _, chunk := range chunks {
		encodedData, err := rsa.EncryptPKCS1v15(rand.Reader, e.pk, chunk)
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(encodedData)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func dataToChunks(data []byte, limit int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(data)/limit+1)
	for len(data) >= limit {
		chunk, data = data[:limit], data[limit:]
		chunks = append(chunks, chunk)
	}
	if len(data) > 0 {
		chunks = append(chunks, data[:])
	}
	return chunks
}
