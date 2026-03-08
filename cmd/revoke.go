package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Revoke(urlStr string) {
	// url can have hash fragment, strip it
	parts := strings.Split(urlStr, "#")
	pasteURL := parts[0]

	req, _ := http.NewRequest("DELETE", pasteURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m failed to send request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m Link not found or already deleted\n")
		os.Exit(1)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		fmt.Fprintf(os.Stderr, "\033[31m ✗ \033[0m unexpected status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("\033[32m ✓ \033[0m Remote paste deleted")
}
