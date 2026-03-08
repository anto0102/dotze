package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"os"
)

type Vault struct {
	Version int               `json:"version"`
	Secrets map[string]string `json:"secrets"`
}

type EncryptedPayload struct {
	IV   string `json:"iv"`
	Data string `json:"data"`
}

func ReadVault(filename string, key []byte) (*Vault, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &Vault{Version: 1, Secrets: make(map[string]string)}, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(content, &payload); err != nil {
		return nil, err
	}

	// Actually, the request says: {"iv":"<base64>","data":"<base64>"}
	// But the Encrypt function returns concatenated nonce+ciphertext.
	// Let's adapt.

	decodedIV, err := base64.StdEncoding.DecodeString(payload.IV)
	if err != nil {
		return nil, err
	}
	decodedData, err := base64.StdEncoding.DecodeString(payload.Data)
	if err != nil {
		return nil, err
	}

	decrypted, err := Decrypt(append(decodedIV, decodedData...), key)
	if err != nil {
		return nil, err
	}

	var vault Vault
	if err := json.Unmarshal(decrypted, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}

func WriteVault(filename string, vault *Vault, key []byte) error {
	plaintext, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	// Encrypt returns nonce+ciphertext
	// nonceSize is usually 12 for GCM
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()

	payload := EncryptedPayload{
		IV:   base64.StdEncoding.EncodeToString(encrypted[:nonceSize]),
		Data: base64.StdEncoding.EncodeToString(encrypted[nonceSize:]),
	}

	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, content, 0600)
}
