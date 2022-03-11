package azcat

import (
	"context"
	"encoding/json"
	"os"
	"text/template"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	format  = pflag.StringP("format", "f", `{{ .Name }}\n`, "Format for list mode")
	formats = map[string]string{
		"long": `{{ .Name | printf '%40s' }} {{ .Size | humanize_size }}\n`,
	}
)

//type renderer func(*azblob.BlobItemInternal)
type blobInfo struct {
	Name         string    `json:"name"`
	Size         uint64    `json:"size"`
	LastModified time.Time `json:"last_modified"`
}
type renderFunc func(*blobInfo)

func List(containerName, prefix string, client *azblob.ServiceClient) {
	containerClient := client.NewContainerClient(containerName)

	pager := containerClient.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{
		Prefix: &prefix,
	})
	if pager == nil {
		log.Fatalf("Cannot list blobs, ListBlobsFlat returned nil")
		return
	}
	if pager.Err() != nil {
		log.WithError(pager.Err()).Fatalf("Cannot list blobs")
		return
	}

	var render renderFunc
	if *format == "json" {
		render = renderJSON()
	} else {
		render = renderTemplate(*format)
	}

	for pager.NextPage(context.TODO()) {
		page := pager.PageResponse()
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

	t := template.New("output")
	t.Parse(templateString)
	t.Funcs(template.FuncMap{
		"humanize_bytes": func(size uint64) string {
			return humanize.Bytes(size)
		},
	})

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
