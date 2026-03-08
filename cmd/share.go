package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}

	vault, err := internal.ReadVault(".dotze", key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to read vault: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "✗ Failed to marshal vault: %v\n", err)
		os.Exit(1)
	}

	// Generate random password
	passwordBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, passwordBytes); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to generate password: %v\n", err)
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

	// Upload to paste.rs
	pasteURL, err := uploadToPaste(b64Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Upload failed: %v\n", err)
		os.Exit(1)
	}

	// Extract ID (e.g., https://paste.rs/abc -> abc)
	parts := strings.Split(strings.TrimRight(pasteURL, "/"), "/")
	id := parts[len(parts)-1]

	fmt.Printf("✓ Share link (expires when pulled):\n")
	fmt.Printf("  https://paste.rs/%s#%s\n\n", id, password)
	fmt.Printf("⚠  Send this link securely. Delete manually with:\n")
	fmt.Printf("  dotze revoke https://paste.rs/%s#%s\n", id, password)
}

func shareLocal(vault *internal.Vault, out string) {
	pass, err := internal.ReadPassword("Enter password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error reading password: %v\n", err)
		os.Exit(1)
	}
	confirm, _ := internal.ReadPassword("Confirm password: ")
	if pass != confirm {
		fmt.Fprintf(os.Stderr, "✗ Passwords do not match\n")
		os.Exit(1)
	}

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Salt generation failed: %v\n", err)
		os.Exit(1)
	}

	key, err := internal.DeriveKey(pass, salt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Key derivation failed: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "✗ Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Exported to %s\n", out)
	fmt.Printf("ℹ  Share this file and your password separately.\n")
	fmt.Printf("ℹ  Recipient runs: dotze pull %s\n", out)
}

func uploadToPaste(data string) (string, error) {
	req, err := http.NewRequest("POST", "https://paste.rs/", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "dotze/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 && resp.StatusCode != 206 && resp.StatusCode != 200 {
		return "", fmt.Errorf("upload failed: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}
