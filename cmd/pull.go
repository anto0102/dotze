package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/anto0102/dotze/internal"
)

func Pull(target string) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		pullOnline(target)
	} else {
		pullLocal(target)
	}
}

func pullOnline(urlStr string) {
	// url format: https://0x0.st/abc.txt#password
	parts := strings.Split(urlStr, "#")
	if len(parts) < 2 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m URL missing password fragment (#)\n")
		os.Exit(1)
	}
	fileURL := parts[0]
	passwordStr := parts[1]

	// Download
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Download failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Remote file not found or expired\n")
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m unexpected status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	body, _ := io.ReadAll(resp.Body)
	b64Data := string(body)

	encrypted, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Invalid base64 data\n")
		os.Exit(1)
	}

	passwordBytes, err := base64.RawURLEncoding.DecodeString(passwordStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Invalid password format\n")
		os.Exit(1)
	}

	decrypted, err := internal.Decrypt(encrypted, passwordBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Invalid password or corrupted data\n")
		os.Exit(1)
	}

	var vault internal.Vault
	if err := json.Unmarshal(decrypted, &vault); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to parse secrets: %v\n", err)
		os.Exit(1)
	}

	if confirmAndImport(&vault) {
		fmt.Printf("\033[34m ℹ \033[0m Remote file will expire automatically.\n")
	}
}

func pullLocal(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to read file: %v\n", err)
		os.Exit(1)
	}

	var exp ExportedVault
	if err := json.Unmarshal(content, &exp); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Invalid file format\n")
		os.Exit(1)
	}

	pass, err := internal.ReadPassword("Enter password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Error: %v\n", err)
		os.Exit(1)
	}

	salt, _ := base64.StdEncoding.DecodeString(exp.Salt)
	key, err := internal.DeriveKey(pass, salt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Key derivation failed\n")
		os.Exit(1)
	}

	iv, _ := base64.StdEncoding.DecodeString(exp.IV)
	data, _ := base64.StdEncoding.DecodeString(exp.Data)

	decrypted, err := internal.Decrypt(append(iv, data...), key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Decryption failed - wrong password?\n")
		os.Exit(1)
	}

	var vault internal.Vault
	if err := json.Unmarshal(decrypted, &vault); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to parse secrets\n")
		os.Exit(1)
	}

	if confirmAndImport(&vault) {
		fmt.Printf("\033[32m ✓ \033[0m Imported from %s\n", filename)
	}
}

func confirmAndImport(vault *internal.Vault) bool {
	keys := make([]string, 0, len(vault.Secrets))
	for k := range vault.Secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Printf("Found %d secrets: %s\n", len(keys), strings.Join(keys, ", "))
	fmt.Print("Import all? [y/N]: ")
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		fmt.Println("\033[31m ✗ \033[0m Import cancelled")
		return false
	}

	// Load local vault
	key, err := internal.GetKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m %v\n", err)
		os.Exit(1)
	}

	localVault, err := internal.ReadVault(".dotze", key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m %v\n", err)
		os.Exit(1)
	}

	// Merge
	for k, v := range vault.Secrets {
		localVault.Secrets[k] = v
	}

	if err := internal.WriteVault(".dotze", localVault, key); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Failed to save vault: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓ \033[0m Imported %d secrets\n", len(vault.Secrets))
	return true
}

func deletePaste(pasteURL string) error {
	req, _ := http.NewRequest("DELETE", pasteURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("delete failed: %d", resp.StatusCode)
	}
	return nil
}
