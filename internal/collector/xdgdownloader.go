package collector

import (
	"os/exec"

	"github.com/pkg/errors"
)

// XDGDownloader is a type that implements the Downloader interface to open a downloadable using xdg-open.
type XDGDownloader struct{}

// Download opens the URI of the downloadable using xdg-open.
func (d XDGDownloader) Download(dl Downloadable) error {
	cmd := exec.Command("xdg-open", dl.URI())

	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "failed to start xdg-open")
	}

	return nil
}
