package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Sign(src []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	return h.Sum(nil)
}
