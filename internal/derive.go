package internal

import (
	"crypto/sha256"

	"crypto/hkdf"
)

func DeriveKey(password string, salt []byte) ([]byte, error) {
	return hkdf.Key(sha256.New, []byte(password), salt, "dotze-share-v1", 32)
}
