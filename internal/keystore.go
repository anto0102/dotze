package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GetKeyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(absPath))
	hashStr := fmt.Sprintf("%x", hash)

	keyDir := filepath.Join(home, ".dotze", "keys")
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(keyDir, hashStr+".key"), nil
}

func LoadOrGenerateKey() ([]byte, error) {
	path, err := GetKeyPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		key := make([]byte, 32) // AES-256
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(key)
		if err := os.WriteFile(path, []byte(encoded), 0600); err != nil {
			return nil, err
		}
		return key, nil
	}

	encoded, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(encoded))
}

func GetKey() ([]byte, error) {
	path, err := GetKeyPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("key not found. run 'dotze init' first")
	}

	encoded, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(encoded))
}
