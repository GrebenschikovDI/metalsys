package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func Sign(src []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(src)
	return fmt.Sprintf("%x", h.Sum(nil))
}
