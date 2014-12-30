package main

import (
	"log"
	"os"

	"code.revolvingcow.com/revolvingcow/go-code/cmd"
)

// Main entry point to the application
func main() {
	err := cmd.NewApp().Run()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
