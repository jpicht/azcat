package internal

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/auth"
	"github.com/sirupsen/logrus"
)

func GetClient(parsed azblob.BlobURLParts) *azblob.ServiceClient {
	clientBuilder := auth.AuthFromEnv()
	if clientBuilder == nil {
		log.Fatal("No client credentials could be detected")
		return nil
	}

	service := (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String()

	log.WithFields(logrus.Fields{
		"service_url": service,
		"container":   parsed.ContainerName,
		"blob":        parsed.BlobName,
	}).Debug()

	serviceClient, err := clientBuilder.CreateClient(service)
	if err != nil {
		log.WithError(err).Fatalf("Cannot create service client")
		return nil
	}

	return &serviceClient
}
