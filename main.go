package main

import (
	"dockVault/helpers"
	"dockVault/internal/pkg/output"
	"dockVault/storage"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func getClient(cfg helpers.Config) (storage.Storage, error) {
	d, err := helpers.NewDocker()
	if err != nil {
		log.Fatal("Failed to start Docker... is it running?")
		return nil, err
	}

	if cfg.AWS != nil {
		s3, err := storage.NewS3(cfg, d)
		if err != nil {
			log.Fatalln("Failed to create s3 connection")
			return nil, err
		}

		return s3, nil
	} else if cfg.Azure != nil {
		az, err := storage.NewAzureWithAD(cfg, d)
		if err != nil {
			log.Fatalln("Failed to create azure connection")
			return nil, err
		}

		return az, nil
	}

	return nil, errors.New("failed to create client")
}

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
		azureCmd.StringVar(&azureAccount, "a", "", "azure storage account")
		azureCmd.StringVar(&azureAccount, "account", "", "azure storage account")
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
				StorageAccount: azureAccount,
				Container:      azureContainer,
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
		var blobName string
		uploadCmd.StringVar(&blobName, "n", "", "name of the file to be saved")
		uploadCmd.StringVar(&blobName, "name", "", "name of the file to saved")

		uploadCmd.Parse(args[2:])
		if len(uploadCmd.Args()) < 1 {
			fmt.Println("usage: dockvault upload <image id | name:tag>")
			return
		}

		imageId := uploadCmd.Arg(0)

		// Get config and do stuff
		cfg, err := helpers.GetConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		strg, err := getClient(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := strg.Upload(storage.UploadParams{ImageId: imageId, BlobName: blobName}); err != nil {
			fmt.Println(err)
			return
		}

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
