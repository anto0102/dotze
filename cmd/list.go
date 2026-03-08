package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/anto0102/dotze/internal"
)

func List(showValues bool) {
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

	keys := make([]string, 0, len(vault.Secrets))
	for k := range vault.Secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if showValues {
			fmt.Printf("%s=%s\n", k, vault.Secrets[k])
		} else {
			fmt.Println(k)
		}
	}
}
