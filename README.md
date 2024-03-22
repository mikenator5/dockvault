# DockVault

DockVault is a tool for storing Docker images in cloud storage. It supports AWS S3 and Azure Blob storage.

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/mikenator5/dockvault.git
cd dockvault
go install .
```

## Usage

Initialize DockVault with your cloud storage provider:

AWS:

```bash
dockvault configure aws --bucket mybucket --region us-west-2
```

Azure:

```bash
dockvault configure az --account myaccount --container mycontainer
```

Upload an image to the vault:

```bash
dockvault upload --name myImageName <image id | name:tag>
```

Load an image from the vault:

```bash
dockvault upload --name myImageName <image id | name:tag>
```

List all images in the vault:

```bash
dockvault list
```
