package main

import (
	"context"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/pkg/azcat"
)

func getExplicitMode() azcat.Mode {
	modes := []azcat.Mode{}

	if *read {
		modes = append(modes, azcat.EMode.List())
	}

	if *list {
		modes = append(modes, azcat.EMode.List())
	}

	if *remove {
		modes = append(modes, azcat.EMode.Remove())
	}

	if *write {
		modes = append(modes, azcat.EMode.Write())
	}

	if len(modes) == 0 {
		return azcat.EMode.None()
	}

	if len(modes) == 1 {
		return modes[0]
	}

	log.Fatal("list, remove and write are mutually exclusive")

	panic("unreachable")
}

func guessMode(blobUrl azblob.BlobURLParts, serviceClient *azblob.ServiceClient) azcat.Mode {
	if blobUrl.BlobName == "/" || blobUrl.BlobName == "" {
		return azcat.EMode.List()
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
		return azcat.EMode.Write()
	}

	if !stdInIsPipe {
		return azcat.EMode.Read()
	}

	return azcat.EMode.None()
}

func isPipe(f *os.File) bool {
	s, err := f.Stat()
	if err != nil {
		log.WithError(err).Fatal("cannot stat STDIN")
	}

	return s.Mode()&os.ModeCharDevice == 0
}
