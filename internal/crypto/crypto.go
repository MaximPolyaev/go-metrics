// Package crypto - пакет для работы с шифрованием
package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyBlock, err := loadBlock(path)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyBlock, err := loadBlock(path)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
}

func loadBlock(path string) (*pem.Block, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(fileData)
	return pemBlock, nil
}
