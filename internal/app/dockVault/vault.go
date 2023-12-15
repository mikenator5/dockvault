package dockVault

import "fmt"

func NewVault(args []string) {
	switch args[0] {
	case "upload":
		if len(args) < 2 {
			printUsage()
			return
		}
		upload(args[1])
	case "list":
		fmt.Println("listing saved images...")
	case "load":
		fmt.Println("loading image...")
	default:
		printUsage()
	}
}

func upload(imageId string) {
	if len(imageId) < 1 {
		printUsage()
		return
	}
	fmt.Printf("Uploading %s\n", imageId)
}

func list() {

}

func load() {

}

func printUsage() {
	fmt.Println("Usage: dockVault <upload | list | load>")
}
