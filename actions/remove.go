package actions

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sirupsen/logrus"
)

func Remove(ctx context.Context, containerName, blobName string, client *azblob.ServiceClient) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("remove")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlobClient(blobName)
	_, err := blobClient.Delete(ctx, nil)
	if err != nil {
		log.WithError(err).Fatal("Delete failed")
	}
}
