package main

import (
	"fmt"
	"os"

	"qed/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
