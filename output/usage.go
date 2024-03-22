package output

import (
	"flag"
	"fmt"
)

// PrintUsage prints the usage information for the dockvault tool.
func PrintUsage() {
	fmt.Println(`dockvault: A tool used for managing Docker images in cloud storage.
Usage:
	dockvault command [arguments]
Commands:
	configure	configure dockvault with your cloud storage provider
	list		list images currently being stored
	load		load an image from the cloud into local docker
	upload		upload in image from Docker to the cloud
	`)
}

// ConfigureUsage prints the usage information for the dockvault configure command.
func ConfigureUsage() {
	w := flag.CommandLine.Output()

	fmt.Fprintf(w, `usage of dockvault configure:
    aws      configure dockvault for aws s3 storage
    az    configure dockvault for azure blob storage
`)
}

// AWSCmdUsage prints the usage information for the "dockvault configure aws" command.
func AWSCmdUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, `usage of dockvault configure aws:
--bucket, -b: name of the AWS bucket
--region, -r: AWS region

Example:
dockvault configure aws --bucket mybucket --region us-west-2
`)
}

// AzureCmdUsage prints the usage information for the dockvault configure az command.
func AzureCmdUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, `usage of dockvault configure az:
    --account, -a:   name of the Azure storage account
    --container, -c: name of the Azure blob container

    Example:
    dockvault configure az --account myaccount --container mycontainer
`)
}

// UploadCmdUsage prints the usage information for the "dockvault upload" command.
func UploadCmdUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, `usage: dockvault upload <image id | name:tag>
    --name, -n: name of the file to be saved. Defaults to image id if none provided

    Example:
    dockvault upload --name myImageName myImage
`)
}

// LoadCmdUsage prints the usage information for the "dockvault load" command.
func LoadCmdUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, `usage: dockvault load <name of object in cloud>

    Example:
    dockvault upload myImage
`)
}
