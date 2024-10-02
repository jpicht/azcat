package internal

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/sirupsen/logrus"
)

func GetClient(parsed azblob.URLParts) *service.Client {
	log := GetLog("internal.GetClient")

	service := (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String()

	log.WithFields(logrus.Fields{
		"service_url": service,
		"container":   parsed.ContainerName,
		"blob":        parsed.BlobName,
	}).Debug()

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.WithError(err).Debug("default creds failed")
		return nil
	}

	azclient, err := azblob.NewClient(service, credential, &azblob.ClientOptions{})
	if err != nil {
		log.WithError(err).Fatalf("Cannot create service client")
		return nil
	}

	return azclient.ServiceClient()
}
