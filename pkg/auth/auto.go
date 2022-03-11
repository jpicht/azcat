package auth

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type (
	ClientCreatorFn func(serviceUrl string) (azblob.ServiceClient, error)
	ClientCreator   struct {
		createClient ClientCreatorFn
	}
)

func (c *ClientCreator) CreateClient(serviceUrl string) (azblob.ServiceClient, error) {
	return c.createClient(serviceUrl)
}

func AuthFromEnv() *ClientCreator {
	log.Debug("trying to create a client from the environment")

	if c := AuthFromEnvConnectionString(); c != nil {
		log.Debug("got client with connection string")
		return c
	}

	if c := AuthFromEnvSharedKey(); c != nil {
		log.Debug("got client with shared key")
		return c
	}

	if c := AuthFromMetadata(); c != nil {
		log.Debug("got client via metadata service")
		return c
	}

	log.Info("no client found")
	return nil
}

func AuthFromEnvSharedKey() *ClientCreator {
	log.Debug("trying AZURE_STORAGE_ACCOUNT_NAME + AZURE_STORAGE_ACCOUNT_KEY")
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		return nil
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		return nil
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil
	}

	return &ClientCreator{
		func(serviceUrl string) (azblob.ServiceClient, error) {
			return azblob.NewServiceClientWithSharedKey(serviceUrl, credential, &azblob.ClientOptions{})
		},
	}
}

func AuthFromEnvConnectionString() *ClientCreator {
	log.Debug("trying AZURE_STORAGE_CONNECTION_STRING")

	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil
	}

	return &ClientCreator{
		func(serviceUrl string) (azblob.ServiceClient, error) {
			return azblob.NewServiceClientFromConnectionString(connectionString, &azblob.ClientOptions{})
		},
	}
}

func AuthFromMetadata() *ClientCreator {
	log.Debug("trying to get identity via metadata service")

	credential, err := azidentity.NewEnvironmentCredential(&azidentity.EnvironmentCredentialOptions{})

	if err != nil {
		log.WithError(err).Debug("metadata failed")
		return nil
	}

	return &ClientCreator{
		func(serviceUrl string) (azblob.ServiceClient, error) {
			return azblob.NewServiceClient(serviceUrl, credential, &azblob.ClientOptions{})
		},
	}
}
