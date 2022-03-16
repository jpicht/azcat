package actions

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func Ping(ctx context.Context, client *azblob.ServiceClient) {
	_, err := client.GetAccountInfo(ctx)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
