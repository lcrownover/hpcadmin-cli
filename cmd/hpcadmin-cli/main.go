package main

import (
	"fmt"
	"os"

	"github.com/lcrownover/hpcadmin-cli/internal/cli"
)

func main() {
	err := cli.Execute()
	if err != nil {
		fmt.Printf("Error executing cli: %v\n", err)
		os.Exit(1)
	}
}
