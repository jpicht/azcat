package actions

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

func Run(mode Mode, bloburl azblob.BlobURLParts, client *azblob.ServiceClient) {
	switch mode {
	case EMode.List():
		List(bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Read():
		Read(bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Write():
		Write(bloburl.ContainerName, bloburl.BlobName, client)
	case EMode.Remove():
		Remove(bloburl.ContainerName, bloburl.BlobName, client)
	default:
		log.WithField("mode", mode.String()).Fatal("Not implemented")
	}
}
