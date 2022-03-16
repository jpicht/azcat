package actions

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func Run(mode Mode, bloburl azblob.BlobURLParts, client *azblob.ServiceClient) {
	context, cancel := context.WithTimeout(context.Background(), getTimeout(mode))
	defer cancel()

	switch mode {
	case EMode.List():
		List(context, bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Read():
		Read(context, bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Write():
		Write(context, bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Remove():
		Remove(context, bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Ping():
		Ping(context, client)
	default:
		log.WithField("mode", mode.String()).Fatal("Not implemented")
	}
}
