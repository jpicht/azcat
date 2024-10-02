package actions

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

func Ping(ctx context.Context, client *service.Client) {
	// we only get here, if client setup was successful
	os.Exit(0)
}
