package main

import (
	"log"
	"os"
	"path/filepath"

	"gcu/internal/cmd"
)

func main() {
	baseName := filepath.Base(os.Args[0])

	err := cmd.NewRootCommand(baseName).Execute()
	if err != nil {
		log.Fatalf("command execution fail: %s", err.Error())
	}
}
