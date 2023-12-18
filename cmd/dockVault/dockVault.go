package main

import (
	"dockVault/internal/app/dockVault"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Unable to load environment variables. Do you have a .env file?")
		return
	}
	dockVault.NewVault()
}
