//go:build windows

package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	fmt.Println()
	return strings.TrimRight(password, "\r\n"), err
}
