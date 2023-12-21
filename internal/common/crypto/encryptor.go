package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var MetEncryptor *Encryptor

type Encryptor struct {
	publicKey *rsa.PublicKey
}

func InitEncryptor(pathFile string) error {

	b, err := os.ReadFile(pathFile)
	if err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	block, _ := pem.Decode(b)
	if block == nil || block.Type != "PUBLIC KEY" {
		return fmt.Errorf("error invalid block type")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error while parsing public key: %w", err)
	}

	MetEncryptor = &Encryptor{
		publicKey: pubKey,
	}

	return nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("error while encrypting data %w", err)
	}

	return ciphertext, nil
}
