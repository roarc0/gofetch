package collector

import (
	"os/exec"

	"github.com/pkg/errors"
)

// WebTorrentDownloader is a type that implements the Downloader interface to open a downloadable using WebTorrent-open.
type WebTorrentDownloader struct{}

// Download opens the URI of the downloadable using WebTorrent-open.
func (d WebTorrentDownloader) Download(dl Downloadable) error {
	cmd := exec.Command("webtorrent", "--mpv", "download", dl.URI())

	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "failed to start WebTorrent-open")
	}

	go func() {
		err = cmd.Wait()
		// log.Info().Err(err).Msg("Webtorrent finished")
	}()

	return nil
}
