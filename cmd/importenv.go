package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tuonome/dotze/internal"
)

func Import(filename string) {
	if filename == "" {
		filename = ".env"
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to open %s: %v\033[0m\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

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

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vault.Secrets[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			count++
		}
	}

	if err := internal.WriteVault(".dotze", vault, key); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Failed to write vault: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m ✓\033[0m Imported %d secrets from %s\n", count, filename)
}
