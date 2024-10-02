package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/jpicht/azcat/actions"
	"github.com/jpicht/azcat/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	log logrus.FieldLogger

	list   = pflag.BoolP("list", "l", false, "Enable list mode")
	ping   = pflag.BoolP("ping", "p", false, "Enable ping mode")
	read   = pflag.BoolP("read", "r", false, "Enable read mode")
	remove = pflag.BoolP("remove", "x", false, "Enable remove mode")
	write  = pflag.BoolP("write", "w", false, "Enable write mode")
)

func main() {
	pflag.Parse()
	log = internal.GetLog("azblob")

	mode := getExplicitMode()

	log.WithField("mode", mode).Debug()

	raw := pflag.Arg(0)

	if raw == "" {
		log.Fatal("Azure URL cannot be empty")
	}

	parsed, err := sas.ParseURL(raw)

	if err != nil {
		log.WithError(err).Fatalf("Invalid URL %#v", raw)
		return
	}

	serviceClient := internal.GetClient(parsed)

	if mode == actions.EMode.None() {
		mode = guessMode(parsed, serviceClient)
		log.WithField("mode", mode).Debug("guessMode")
	}

	actions.Run(mode, parsed, serviceClient)
}
