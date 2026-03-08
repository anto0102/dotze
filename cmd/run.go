package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/anto0102/dotze/internal"
)

func Run(args []string) {
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

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ No command provided. Use: dotze run -- <command>\033[0m\n")
		os.Exit(1)
	}

	env := os.Environ()
	for k, v := range vault.Secrets {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Command failed: %v\033[0m\n", err)
		os.Exit(1)
	}
}
