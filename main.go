package main

import (
	"fmt"
	"os"

	"github.com/anto0102/dotze/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "init":
		cmd.Init()
	case "set":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "\033[31m ✗ Usage: dotze set <KEY> <VALUE>\033[0m\n")
			os.Exit(1)
		}
		cmd.Set(os.Args[2], os.Args[3])
	case "get":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "\033[31m ✗ Usage: dotze get <KEY>\033[0m\n")
			os.Exit(1)
		}
		cmd.Get(os.Args[2])
	case "list":
		showValues := false
		if len(os.Args) > 2 && os.Args[2] == "--show-values" {
			showValues = true
		}
		cmd.List(showValues)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "\033[31m ✗ Usage: dotze delete <KEY>\033[0m\n")
			os.Exit(1)
		}
		cmd.Delete(os.Args[2])
	case "run":
		// Find "--" separator
		var runArgs []string
		found := false
		for i, arg := range os.Args {
			if arg == "--" {
				runArgs = os.Args[i+1:]
				found = true
				break
			}
		}
		if !found {
			if len(os.Args) > 2 {
				runArgs = os.Args[2:]
			} else {
				fmt.Fprintf(os.Stderr, "\033[31m ✗ Usage: dotze run -- <command>\033[0m\n")
				os.Exit(1)
			}
		}
		cmd.Run(runArgs)
	case "import":
		filename := ""
		if len(os.Args) > 2 {
			filename = os.Args[2]
		}
		cmd.Import(filename)
	case "share":
		local := false
		out := "secrets.enc"
		for i := 2; i < len(os.Args); i++ {
			arg := os.Args[i]
			if arg == "--local" {
				local = true
			}
			if arg == "--out" && i+1 < len(os.Args) {
				out = os.Args[i+1]
				i++
			}
		}
		cmd.Share(local, out)
	case "pull":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "✗ Usage: dotze pull <url-or-file>\n")
			os.Exit(1)
		}
		cmd.Pull(os.Args[2])
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "\033[31m ✗ Unknown command: %s\033[0m\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("dotze — Lightweight CLI secret manager")
	fmt.Println("\nUsage:")
	fmt.Println("  dotze init                   Initialize a new vault")
	fmt.Println("  dotze set <KEY> <VALUE>      Set a secret")
	fmt.Println("  dotze get <KEY>              Get a secret value")
	fmt.Println("  dotze list [--show-values]   List all secrets")
	fmt.Println("  dotze delete <KEY>           Delete a secret")
	fmt.Println("  dotze run -- <command>       Run a command with secrets as env vars")
	fmt.Println("  dotze import [file]          Import secrets from .env file (default: .env)")
	fmt.Println("  dotze share                  Share secrets via encrypted link")
	fmt.Println("  dotze share --local          Export secrets to encrypted file")
	fmt.Println("  dotze pull <url|file>        Import secrets from link or file")
	fmt.Println("  dotze help                   Show this help")
}
