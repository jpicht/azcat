package actions

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sirupsen/logrus"
)

func Write(ctx context.Context, containerName, blobName string, client *azblob.ServiceClient) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("write")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobName)

	_, err := blobClient.UploadStreamToBlockBlob(ctx, os.Stdin, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		log.WithError(err).Fatal("Write failed")
	}
}
