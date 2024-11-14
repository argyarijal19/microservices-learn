package lib

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAPIKey membuat kunci acak sepanjang 32 byte dan mengembalikannya dalam format hex
func GenerateAPIKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return ""
	}

	return hex.EncodeToString(key)
}
