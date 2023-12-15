package main

import (
	"dockVault/internal/app/dockVault"
	"os"
)

func main() {
	args := os.Args[1:]
	dockVault.NewVault(args)
}
