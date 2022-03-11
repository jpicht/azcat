package main

import (
	"context"
	"net/url"
	"os"
	"strings"

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
	read   = pflag.BoolP("read", "r", false, "Enable read mode")
	remove = pflag.BoolP("remove", "x", false, "Enable remove mode")
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

	if mode == azcat.EMode.None() {
		mode = guessMode(parsed, serviceClient)
	}

	azcat.Run(mode, parsed, serviceClient)
}

func getMode() azcat.Mode {
	modes := []azcat.Mode{}

	if *read {
		modes = append(modes, azcat.EMode.List())
	}

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
		return azcat.EMode.None()
	}

	if len(modes) == 1 {
		return modes[0]
	}

	log.Fatal("list, remove and write are mutually exclusive")

	panic("unreachable")
}

func guessMode(blobUrl azblob.BlobURLParts, serviceClient *azblob.ServiceClient) azcat.Mode {
	if blobUrl.BlobName == "/" || blobUrl.BlobName == "" {
		return azcat.EMode.List()
	}

	blobClient := serviceClient.NewContainerClient(blobUrl.ContainerName).NewBlobClient(blobUrl.BlobName)
	_, err := blobClient.GetProperties(context.TODO(), &azblob.GetBlobPropertiesOptions{})

	var (
		exists      = true
		stdInIsPipe = isPipe(os.Stdin)
	)

	if err != nil {
		if strings.Contains(err.Error(), "BlobNotFound") {
			exists = false
		} else {
			log.WithError(err).Fatal("GetProperties request failed")
		}
	}

	if !exists {
		return azcat.EMode.Write()
	}

	if !stdInIsPipe {
		return azcat.EMode.Read()
	}

	return azcat.EMode.None()
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

func isPipe(f *os.File) bool {
	s, err := f.Stat()
	if err != nil {
		log.WithError(err).Fatal("cannot stat STDIN")
	}

	return s.Mode()&os.ModeCharDevice == 0
}
