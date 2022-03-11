package main

import (
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/jpicht/azcat/pkg/auth"
	"github.com/jpicht/azcat/pkg/azcat"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	log = logrus.StandardLogger().WithField("module", "main")

	debug  = pflag.BoolP("debug", "d", false, "Enable debug output")
	list   = pflag.BoolP("list", "l", false, "Enable list mode")
	remove = pflag.BoolP("remove", "r", false, "Enable remove mode")
	write  = pflag.BoolP("write", "w", false, "Enable write mode")
)

func main() {
	pflag.Parse()

	if *debug {
		logrus.StandardLogger().SetLevel(logrus.DebugLevel)
	}

	mode := getMode()

	parsed, err := azblob.NewBlobURLParts(pflag.Arg(0))

	if err != nil {
		log.WithError(err).Fatalf("Invalid URL %#v", pflag.Arg(0))
		return
	}

	serviceClient := getClient(parsed)

	azcat.Run(mode, parsed, serviceClient)
}

func getMode() azcat.Mode {
	modes := []azcat.Mode{}

	if *list {
		modes = append(modes, azcat.EMode.List())
	}

	if *remove {
		modes = append(modes, azcat.EMode.Remove())
	}

	if *write {
		modes = append(modes, azcat.EMode.Write())
	}

	if len(modes) == 0 {
		return azcat.EMode.Read()
	}

	if len(modes) == 1 {
		return modes[0]
	}

	log.Fatal("list, remove and write are mutually exclusive")

	panic("unreachable")
}

func getClient(parsed azblob.BlobURLParts) *azblob.ServiceClient {
	clientBuilder := auth.AuthFromEnv()
	if clientBuilder == nil {
		log.Fatal("No client credentials could be detected")
		return nil
	}

	service := (&url.URL{Scheme: parsed.Scheme, Host: parsed.Host}).String()

	log.WithFields(logrus.Fields{
		"service_url": service,
		"container":   parsed.ContainerName,
		"blob":        parsed.BlobName,
	}).Debug()

	serviceClient, err := clientBuilder.CreateClient(service)
	if err != nil {
		log.WithError(err).Fatalf("Cannot create service client")
		return nil
	}

	return &serviceClient
}
