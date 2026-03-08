package cmd

import (
	"fmt"
	"os"

	"github.com/tuonome/dotze/internal"
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

	for k, v := range vault.Secrets {
		if showValues {
			fmt.Printf("%s=%s\n", k, v)
		} else {
			fmt.Println(k)
		}
	}
}
