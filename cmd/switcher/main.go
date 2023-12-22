package main

import (
	"log"
	"os"
	"switcher"
)

func main() {
	// Working directory for nix - pwd()
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := switcher.NewLogger(workDir, true)
	app := switcher.NewApp(workDir, logger)

	app.Serve()
}
