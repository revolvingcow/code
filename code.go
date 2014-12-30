package main

import (
	"log"

	"code.revolvingcow.com/revolvingcow/go-code/cmd"
)

// Main entry point to the application
func main() {
	err := cmd.NewApp().Run()

	// If there was an error report and exit
	if err != nil {
		log.Fatal(err)
	}
}
