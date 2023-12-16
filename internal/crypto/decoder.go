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
	encodedDataChunks := dataToChunks(data, d.size)
	var decodedDataBuf bytes.Buffer

	for _, encodedDataChunk := range encodedDataChunks {
		decodedDataChunk, err := rsa.DecryptPKCS1v15(rand.Reader, d.pk, encodedDataChunk)
		if err != nil {
			return nil, err
		}
		decodedDataBuf.Write(decodedDataChunk)
	}

	return decodedDataBuf.Bytes(), nil
}
