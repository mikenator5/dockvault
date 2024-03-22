package main

import (
	"dockVault/helpers"
	"dockVault/internal/pkg/output"
	"flag"
	"fmt"
	"os"
)

func main() {

	args := os.Args
	if len(args) <= 1 {
		output.PrintUsage()
		return
	}

	switch args[1] {
	case "configure":
		if len(args) <= 2 {
			fmt.Println("usage: dockvault configure <aws|azure>")
			return
		}
		awsCmd := flag.NewFlagSet("aws", flag.ExitOnError)
		var awsBucket string
		awsCmd.StringVar(&awsBucket, "bucket", "", "name of the bucket")
		awsCmd.StringVar(&awsBucket, "b", "", "name of the bucket")
		var awsRegion string
		awsCmd.StringVar(&awsRegion, "region", "", "aws region")
		awsCmd.StringVar(&awsRegion, "r", "", "aws region")

		// az -a,--account -c,--container
		azureCmd := flag.NewFlagSet("az", flag.ExitOnError)
		var azureAccount string
		azureCmd.StringVar(&azureAccount, "a", "", "azure account")
		azureCmd.StringVar(&azureAccount, "account", "", "azure account")
		var azureContainer string
		azureCmd.StringVar(&azureContainer, "c", "", "azure blob container")
		azureCmd.StringVar(&azureContainer, "container", "", "azure blob container")

		switch args[2] {
		case "aws":
			awsCmd.Parse(args[3:])
			if len(awsBucket) < 1 && len(awsRegion) < 1 {
				awsCmd.Usage()
				return
			}

			aws := helpers.AWS{
				Bucket: awsBucket,
				Region: awsRegion,
			}

			err := helpers.NewAWSConfig(&aws)
			if err != nil {
				fmt.Println("Error creating aws config file")
				return
			}
		case "az":
			azureCmd.Parse(args[3:])
			if len(azureAccount) < 1 && len(azureContainer) < 1 {
				azureCmd.Usage()
				return
			}

			az := helpers.Azure{
				Account:   azureAccount,
				Container: azureContainer,
			}

			if err := helpers.NewAzureConfig(&az); err != nil {
				fmt.Println("Error creating azure config file")
				return
			}
		default:
			fmt.Println("usage: dockvault configure <aws|azure>")
		}

	case "upload":
		uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
		uploadCmd.Parse(args[1:])
	default:
		fmt.Printf("Unknown subcommand '%s'\n", args[1])
	}

	// c, err := helpers.NewConfig()
	// if err != nil {
	// 	log.Fatalln("Config failed", err)
	// }
	// d, err := helpers.NewDocker()
	// if err != nil {
	// 	log.Fatal("Failed to start Docker... is it running?")
	// }

}
