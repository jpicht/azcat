package actions

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sirupsen/logrus"
)

func Write(containerName, blobName string, client *azblob.ServiceClient) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("write")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobName)

	_, err := blobClient.UploadStreamToBlockBlob(context.TODO(), os.Stdin, azblob.UploadStreamToBlockBlobOptions{})
	if err != nil {
		log.WithError(err).Fatal("Write failed")
	}
}
