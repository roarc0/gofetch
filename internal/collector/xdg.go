package collector

import (
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// XDGDownloader is a type that implements the Downloader interface to open a downloadable using xdg-open.
type XDGDownloader struct {
	Downloadable
}

// Download opens the URI of the downloadable using xdg-open.
func (o XDGDownloader) Open() error {
	cmd := exec.Command("xdg-open", o.URI())

	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "failed to start xdg-open")
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrapf(err, "failed to wait for xdg-open")
	}

	go func() {
		<-time.After(1 * time.Second)
		cmd.Process.Release()
	}()

	return nil
}
