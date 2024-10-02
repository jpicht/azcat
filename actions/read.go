package actions

import (
	"context"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/sirupsen/logrus"
)

func Read(ctx context.Context, containerName, blobName string, client *service.Client) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("read")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlobClient(blobName)
	stream, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		log.WithError(err).Fatal("Download failed")
	}

	retryReader := stream.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	defer retryReader.Close()

	_, err = io.Copy(os.Stdout, retryReader)
	if err != nil {
		log.WithError(err).Fatal("Download failed")
	}
}
