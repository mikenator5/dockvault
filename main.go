package main

import (
	"dockVault/helpers"
	"log"
)

func main() {
	_, err := helpers.NewDocker()
	if err != nil {
		log.Fatal("Failed to start Docker... is it running?")
	}
}
