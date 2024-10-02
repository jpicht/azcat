package actions

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/jpicht/azcat/internal"
	"github.com/spf13/pflag"
)

func Main(mode Mode) {
	pflag.Parse()

	log := internal.GetLog("internal.Main")

	if pflag.NArg() != 1 {
		log.Fatal("Invalid number of arguments")
	}

	raw := pflag.Arg(0)

	if raw == "" {
		log.Fatal("Azure URL cannot be empty")
	}

	parsed, err := sas.ParseURL(raw)

	if err != nil {
		log.WithError(err).WithField("url", raw).Fatalf("invalid URL")
	}

	serviceClient := internal.GetClient(parsed)

	Run(mode, parsed, serviceClient)
}
