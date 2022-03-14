package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/actions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func Main(mode actions.Mode) {
	pflag.Parse()

	if *debug {
		logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	}

	if pflag.NArg() != 1 {
		log.Fatal("Invalid number of arguments")
	}

	parsed, err := azblob.NewBlobURLParts(pflag.Arg(0))

	if err != nil {
		log.WithError(err).WithField("url", pflag.Arg(0)).Fatalf("invalid URL")
	}

	serviceClient := GetClient(parsed)

	actions.Run(mode, parsed, serviceClient)
}
