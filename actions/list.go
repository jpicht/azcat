package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	format  = pflag.StringP("format", "f", "default", "Format for list mode")
	formats = map[string]string{
		"default": "{{ .Name }}\n",
		"long":    "{{ printf \"%40s\" .Name }}{{ printf \"%10s\" (humanize_bytes .Size) }}  {{ .LastModified }}\n",
	}
)

// type renderer func(*azblob.BlobItemInternal)
type blobInfo struct {
	Name         string    `json:"name"`
	Size         uint64    `json:"size"`
	LastModified time.Time `json:"last_modified"`
}
type renderFunc func(*blobInfo)

func List(ctx context.Context, containerName, prefix string, client *service.Client) {
	if containerName == "" {
		ListContainers(ctx, client)
	} else {
		ListBlobs(ctx, containerName, prefix, client)
	}
}

func ListContainers(ctx context.Context, client *service.Client) {
	log.Debug("ListContainers")
	pager := client.NewListContainersPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.WithError(err).Fatal("Cannot list containers")
			return
		}
		for _, container := range page.ContainerItems {
			fmt.Println(*container.Name)
		}
	}
}

func ListBlobs(ctx context.Context, containerName, prefix string, client *service.Client) {
	log.WithField("container", containerName).Debug("ListBlobs")

	containerClient := client.NewContainerClient(containerName)

	pager := containerClient.NewListBlobsFlatPager(&azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})
	if pager == nil {
		log.Fatalf("Cannot list blobs, ListBlobsFlat returned nil")
		return
	}

	var render renderFunc
	if *format == "json" {
		render = renderJSON()
	} else {
		render = renderTemplate(*format)
	}

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.WithError(err).Fatalf("Cannot list blobs")
			return
		}
		if page.Segment == nil || len(page.Segment.BlobItems) == 0 {
			continue
		}
		for _, blob := range page.Segment.BlobItems {
			bi := &blobInfo{
				Name:         *blob.Name,
				Size:         uint64(*blob.Properties.ContentLength),
				LastModified: *blob.Properties.LastModified,
			}

			render(bi)
		}
	}
}

func renderTemplate(format string) func(*blobInfo) {
	var templateString string
	if format_, ok := formats[format]; ok {
		templateString = format_
	} else {
		templateString = format
	}

	t := template.Must(
		template.New("output").Funcs(template.FuncMap{
			"humanize_bytes": func(size uint64) string {
				return humanize.Bytes(size)
			},
		}).Parse(templateString),
	)

	return func(bi *blobInfo) {
		err := t.Execute(os.Stdout, bi)
		if err != nil {
			log.WithFields(logrus.Fields{
				"blob":     bi,
				"format":   format,
				"template": templateString,
			}).WithError(err).Fatalf("Cannot format blob")
		}
	}
}

func renderJSON() func(*blobInfo) {
	writer := json.NewEncoder(os.Stdout)
	return func(bi *blobInfo) {
		writer.Encode(bi)
	}
}
