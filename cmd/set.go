package cmd

import (
	"fmt"
	"os"

	"github.com/tuonome/dotze/internal"
)

func Set(keyStr, value string) {
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

	vault.Secrets[keyStr] = value

	if err := internal.WriteVault(".dotze", vault, key); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to write vault: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓\033[0m %s saved\n", keyStr)
}
