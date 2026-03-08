package cmd

import (
	"fmt"
	"os"

	"github.com/anto0102/dotze/internal"
)

func Get(keyStr string) {
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

	val, ok := vault.Secrets[keyStr]
	if !ok {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Key %s not found\033[0m\n", keyStr)
		os.Exit(1)
	}

	fmt.Println(val)
}
