package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/actions"
	"github.com/jpicht/azcat/internal"
	"github.com/spf13/pflag"
)

var (
	log = internal.GetLog("azblob")

	list   = pflag.BoolP("list", "l", false, "Enable list mode")
	ping   = pflag.BoolP("ping", "p", false, "Enable ping mode")
	read   = pflag.BoolP("read", "r", false, "Enable read mode")
	remove = pflag.BoolP("remove", "x", false, "Enable remove mode")
	write  = pflag.BoolP("write", "w", false, "Enable write mode")
)

func main() {
	pflag.Parse()

	mode := getExplicitMode()

	parsed, err := azblob.NewBlobURLParts(pflag.Arg(0))

	if err != nil {
		log.WithError(err).Fatalf("Invalid URL %#v", pflag.Arg(0))
		return
	}

	serviceClient := internal.GetClient(parsed)

	if mode == actions.EMode.None() {
		mode = guessMode(parsed, serviceClient)
	}

	actions.Run(mode, parsed, serviceClient)
}
