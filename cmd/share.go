package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/anto0102/dotze/internal"
)

type ExportedVault struct {
	Version int    `json:"version"`
	Salt    string `json:"salt"`
	IV      string `json:"iv"`
	Data    string `json:"data"`
}

func Share(local bool, out string) {
	key, err := internal.GetKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m %v\n", err)
		os.Exit(1)
	}

	vault, err := internal.ReadVault(".dotze", key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to read vault: %v\n", err)
		os.Exit(1)
	}

	if len(vault.Secrets) == 0 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m No secrets to share. Add some with: dotze set KEY value\n")
		os.Exit(1)
	}

	if local {
		shareLocal(vault, out)
	} else {
		shareOnline(vault)
	}
}

func shareOnline(vault *internal.Vault) {
	plaintext, err := json.Marshal(vault)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to marshal vault: %v\n", err)
		os.Exit(1)
	}

	// Generate random password
	passwordBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, passwordBytes); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to generate password: %v\n", err)
		os.Exit(1)
	}
	password := base64.RawURLEncoding.EncodeToString(passwordBytes)

	// Encrypt with password (AES-256-GCM)
	// We use standard Encrypt logic but with the random password
	encrypted, err := internal.Encrypt(plaintext, passwordBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Encryption failed: %v\n", err)
		os.Exit(1)
	}

	// Convert to base64
	b64Data := base64.StdEncoding.EncodeToString(encrypted)

	// Upload to 0x0.st
	pasteURL, err := uploadToPaste(b64Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Upload failed: %v\n", err)
		os.Exit(1)
	}

	// No parsing needed for 0x0.st as we use the direct URL

	fmt.Printf("\033[32m ✓ \033[0m Share link (expires in 24h):\n")
	fmt.Printf("  %s#%s\n\n", pasteURL, password)
	fmt.Printf("\033[33m ⚠ \033[0m This link cannot be manually revoked.\n")
	fmt.Printf("    It expires automatically after 24 hours.\n")
}

func shareLocal(vault *internal.Vault, out string) {
	pass, err := internal.ReadPassword("Enter password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Error reading password: %v\n", err)
		os.Exit(1)
	}
	confirm, _ := internal.ReadPassword("Confirm password: ")
	if pass != confirm {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Passwords do not match\n")
		os.Exit(1)
	}

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Salt generation failed: %v\n", err)
		os.Exit(1)
	}

	key, err := internal.DeriveKey(pass, salt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Key derivation failed: %v\n", err)
		os.Exit(1)
	}

	plaintext, _ := json.Marshal(vault)
	encrypted, err := internal.Encrypt(plaintext, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Encryption failed: %v\n", err)
		os.Exit(1)
	}

	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()

	exp := ExportedVault{
		Version: 1,
		Salt:    base64.StdEncoding.EncodeToString(salt),
		IV:      base64.StdEncoding.EncodeToString(encrypted[:nonceSize]),
		Data:    base64.StdEncoding.EncodeToString(encrypted[nonceSize:]),
	}

	content, _ := json.Marshal(exp)
	if err := os.WriteFile(out, content, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓ \033[0m Exported to %s\n", out)
	fmt.Printf("\033[34m ℹ \033[0m Share this file and your password separately.\n")
	fmt.Printf("\033[34m ℹ \033[0m Recipient runs: dotze pull %s\n", out)
}

func uploadToPaste(data string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "secrets")
	if err != nil {
		return "", err
	}
	part.Write([]byte(data))
	writer.Close()

	req, err := http.NewRequest("POST", "https://0x0.st", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "dotze/2.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("upload failed: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}
