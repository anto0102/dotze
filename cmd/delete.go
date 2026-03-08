package cmd

import (
	"fmt"
	"os"

	"github.com/anto0102/dotze/internal"
)

func Delete(keyStr string) {
	key, err := internal.GetKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ %v\033[0m\n", err)
		os.Exit(1)
	}

	vault, err := internal.ReadVault(".dotze", key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to read vault: %v\033[0m\n", err)
		os.Exit(1)
	}

	if _, ok := vault.Secrets[keyStr]; !ok {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Key %s not found\033[0m\n", keyStr)
		os.Exit(1)
	}

	delete(vault.Secrets, keyStr)

	if err := internal.WriteVault(".dotze", vault, key); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to write vault: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓\033[0m %s deleted\n", keyStr)
}
