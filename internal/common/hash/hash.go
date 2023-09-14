package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Sign(src, key []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	_, err := h.Write(src)
	if err != nil {
		return nil, err
	}
	dst := h.Sum(nil)
	return dst, nil
}
