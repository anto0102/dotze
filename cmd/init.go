package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/tuonome/dotze/internal"
)

func Init() {
	if _, err := os.Stat(".dotze"); err == nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ .dotze already exists\033[0m\n")
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓\033[0m Generating AES key...\n")
	key, err := internal.LoadOrGenerateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to generate key: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓\033[0m Key saved in local keystore\n")

	vault := &internal.Vault{
		Version: 1,
		Secrets: make(map[string]string),
	}

	if err := internal.WriteVault(".dotze", vault, key); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to create vault: %v\033[0m\n", err)
		os.Exit(1)
	}
	fmt.Printf("\033[32m ✓\033[0m Vault created (.dotze)\n")

	if _, err := os.Stat(".gitignore"); err == nil {
		content, _ := os.ReadFile(".gitignore")
		if !strings.Contains(string(content), ".dotze.key") {
			f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString("\n# dotze key path if local (default is ~/.dotze/keys/)\n.dotze.key\n")
				f.Close()
				fmt.Printf("\033[32m ✓\033[0m Updated .gitignore\n")
			}
		}
	} else {
		os.WriteFile(".gitignore", []byte(".dotze.key\n"), 0644)
		fmt.Printf("\033[32m ✓\033[0m Created .gitignore\n")
	}

	fmt.Printf("\033[32m ✓\033[0m Initialization complete\n")
}
