package main

import (
	"fmt"
	"os"

	"github.com/a630140621/lethe/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
