package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(src []byte, key string) (hash string) {
	secretkey := []byte("key")
	h := hmac.New(sha256.New, secretkey)
	h.Write(src)
	hsh := h.Sum(nil)
	return hex.EncodeToString(hsh)
}
