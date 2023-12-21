package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var MetDecryptor *Decrypt

type Decrypt struct {
	privateKey *rsa.PrivateKey
}

func InitDecryptor(pathFile string) error {
	b, err := os.ReadFile(pathFile)
	if err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	block, _ := pem.Decode(b)
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("error invalid block type")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error while parsing private key: %w", err)
	}

	MetDecryptor = &Decrypt{
		privateKey: privateKey,
	}

	return nil
}

func (d *Decrypt) Decrypt(data []byte) ([]byte, error) {

	decryptedMsg, err := rsa.DecryptPKCS1v15(rand.Reader, d.privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding data %w", err)
	}

	return decryptedMsg, nil
}
