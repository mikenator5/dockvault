package main

import (
	"dockVault/helpers"
	"fmt"
	"log"
)

func main() {
	// c, err := helpers.NewConfig()
	// if err != nil {
	// 	log.Fatalln("Config failed", err)
	// }
	d, err := helpers.NewDocker()
	if err != nil {
		log.Fatal("Failed to start Docker... is it running?")
	}
	fmt.Println(d)
	// s3, err := storage.NewS3(c, d)
	// if err != nil {
	// 	log.Fatalln("S3 Failed", err)
	// }
	// s3.List()
}
