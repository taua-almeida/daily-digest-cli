package main

import (
	"log"
	"os"

	"github.com/taua-almeida/gh-daily-digest-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
