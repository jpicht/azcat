package actions

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/pflag"
)

var (
	timeout = pflag.String("timeout", "1s", "Timeout")
)

func Ping(client *azblob.ServiceClient) {
	duration, err := time.ParseDuration(*timeout)
	if err != nil {
		log.WithError(err).Fatal("Invalid timeout specified")
	}

	ctx, cf := context.WithTimeout(context.TODO(), duration)
	defer cf()

	_, err = client.GetAccountInfo(ctx)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
