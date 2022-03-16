package main

import (
	"context"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/actions"
)

func getExplicitMode() actions.Mode {
	modes := []actions.Mode{}

	if *read {
		modes = append(modes, actions.EMode.List())
	}

	if *list {
		modes = append(modes, actions.EMode.List())
	}

	if *remove {
		modes = append(modes, actions.EMode.Remove())
	}

	if *ping {
		modes = append(modes, actions.EMode.Ping())
	}

	if *write {
		modes = append(modes, actions.EMode.Write())
	}

	if len(modes) == 0 {
		return actions.EMode.None()
	}

	if len(modes) == 1 {
		return modes[0]
	}

	log.Fatal("list, ping, read, remove and write are mutually exclusive")

	panic("unreachable")
}

func guessMode(blobUrl azblob.BlobURLParts, serviceClient *azblob.ServiceClient) actions.Mode {
	switch os.Args[0] {
	case "azls":
		return actions.EMode.List()
	case "azcat":
		return actions.EMode.Read()
	case "azput":
		return actions.EMode.Write()
	case "azping":
		return actions.EMode.Ping()
	case "azrm":
		return actions.EMode.Remove()
	}

	if blobUrl.BlobName == "/" || blobUrl.BlobName == "" {
		return actions.EMode.List()
	}

	blobClient := serviceClient.NewContainerClient(blobUrl.ContainerName).NewBlobClient(blobUrl.BlobName)
	_, err := blobClient.GetProperties(context.TODO(), &azblob.GetBlobPropertiesOptions{})

	var (
		exists      = true
		stdInIsPipe = isPipe(os.Stdin)
	)

	if err != nil {
		if strings.Contains(err.Error(), "BlobNotFound") {
			exists = false
		} else {
			log.WithError(err).Fatal("GetProperties request failed")
		}
	}

	if !exists {
		return actions.EMode.Write()
	}

	if !stdInIsPipe {
		return actions.EMode.Read()
	}

	return actions.EMode.None()
}

func isPipe(f *os.File) bool {
	s, err := f.Stat()
	if err != nil {
		log.WithError(err).Fatal("cannot stat STDIN")
	}

	return s.Mode()&os.ModeCharDevice == 0
}
