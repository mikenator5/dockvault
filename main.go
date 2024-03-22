package main

import (
	"dockvault/helpers"
	"dockvault/output"
	"dockvault/storage"

	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

// getClient returns a storage client based on the provided configuration.
// It creates a Docker client and initializes the appropriate storage client based on the configuration.
// If the configuration specifies AWS, it creates an S3 storage client.
// If the configuration specifies Azure, it creates an Azure storage client with Active Directory authentication.
// If the configuration does not match any supported storage options, it returns an error.
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

	// Handle configuration
	if args[1] == "configure" {
		if len(args) <= 2 {
			output.ConfigureUsage()
			return
		}
		awsCmd := flag.NewFlagSet("aws", flag.ExitOnError)
		awsCmd.Usage = output.AWSCmdUsage
		var awsBucket string
		awsCmd.StringVar(&awsBucket, "bucket", "", "name of the bucket")
		awsCmd.StringVar(&awsBucket, "b", "", "name of the bucket")
		var awsRegion string
		awsCmd.StringVar(&awsRegion, "region", "", "aws region")
		awsCmd.StringVar(&awsRegion, "r", "", "aws region")

		// az -a,--account -c,--container
		azureCmd := flag.NewFlagSet("az", flag.ExitOnError)
		azureCmd.Usage = output.AzureCmdUsage
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
			output.ConfigureUsage()
		}
		return
	}

	// Begin main program
	cfg, err := helpers.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	strgClient, err := getClient(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch args[1] {
	case "upload":
		uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
		uploadCmd.Usage = output.UploadCmdUsage
		var blobName string
		uploadCmd.StringVar(&blobName, "n", "", "name of the file to be saved")
		uploadCmd.StringVar(&blobName, "name", "", "name of the file to saved")

		uploadCmd.Parse(args[2:])
		if len(uploadCmd.Args()) < 1 {
			uploadCmd.Usage()
			return
		}

		imageId := uploadCmd.Arg(0)

		if err := strgClient.Upload(storage.UploadParams{ImageId: imageId, BlobName: blobName}); err != nil {
			fmt.Println(err)
			return
		}
	case "load":
		loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
		loadCmd.Usage = output.LoadCmdUsage
		loadCmd.Parse(args[2:])
		if len(loadCmd.Args()) < 1 {
			loadCmd.Usage()
			return
		}

		blobName := loadCmd.Arg(0)

		if err := strgClient.Load(storage.LoadParams{BlobName: blobName}); err != nil {
			fmt.Println(err)
			return
		}
	case "list":
		if err := strgClient.List(); err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Printf("Unknown subcommand '%s'\n", args[1])
	}
}
