package actions

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sirupsen/logrus"
)

type adapter struct {
	s io.Seeker
	w io.Writer
	n int64
}

func (a *adapter) pad(diff int64) error {
	// try seeking first
	if a.s != nil {
		n, err := a.s.Seek(diff, io.SeekCurrent)
		a.n += n
		if err == nil {
			return nil
		}
	}

	// pad with zero bytes
	buff := make([]byte, 1024)
	var (
		n   int
		err error
	)
	for diff > 0 {
		if diff >= 1024 {
			n, err = a.w.Write(buff)
		} else {
			n, err = a.w.Write(buff[0:diff])
		}

		a.n += int64(n)
		diff -= int64(n)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *adapter) WriteAt(p []byte, off int64) (n int, err error) {
	if a.n > off {
		return 0, fmt.Errorf("illegal seek")
	}
	if a.n < off {
		err = a.pad(off - a.n)
		if err != nil {
			return
		}
	}

	n, err = a.w.Write(p)

	a.n += int64(n)

	return
}

func Read(containerName, blobName string, client *azblob.ServiceClient) {
	log.WithFields(logrus.Fields{
		"container": containerName,
		"blob":      blobName,
	}).Debug("read")
	containerClient := client.NewContainerClient(containerName)
	blobClient := containerClient.NewBlobClient(blobName)
	err := blobClient.DownloadBlobToWriterAt(context.TODO(), 0, azblob.CountToEnd, &adapter{os.Stdout, os.Stdout, 0}, azblob.HighLevelDownloadFromBlobOptions{
		Parallelism: 1,
	})
	if err != nil {
		log.WithError(err).Fatal("Download failed")
	}
}
