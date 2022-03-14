package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/actions"
	"github.com/jpicht/azcat/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	log = logrus.StandardLogger().WithField("module", "main")

	debug  = pflag.BoolP("debug", "d", false, "Enable debug output")
	list   = pflag.BoolP("list", "l", false, "Enable list mode")
	read   = pflag.BoolP("read", "r", false, "Enable read mode")
	remove = pflag.BoolP("remove", "x", false, "Enable remove mode")
	write  = pflag.BoolP("write", "w", false, "Enable write mode")
)

func main() {
	pflag.Parse()

	if *debug {
		logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	}

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
