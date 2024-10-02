package actions

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/sirupsen/logrus"
)

func Write(ctx context.Context, containerName, blobName string, client *service.Client) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("write")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobName)

	_, err := blobClient.UploadStream(ctx, os.Stdin, &blockblob.UploadStreamOptions{})
	if err != nil {
		log.WithError(err).Fatal("Write failed")
	}
}
