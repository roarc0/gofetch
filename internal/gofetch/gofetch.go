package gofetch

import (
	"context"
	"errors"
	"fmt"

	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/config"
	"github.com/roarc0/gofetch/internal/memory"
)

const (
	downloadsBucket = "downloads"
)

type Downloadable struct {
	collector.Downloadable
	Optional bool
	Seen     bool // TODO see if it has been downloaded or ignored
}

type GoFetch struct {
	cfg        *config.Config
	memory     memory.Memory
	downloader collector.Downloader
}

// NewGoFetch creates a new GoFetch object.
func NewGoFetch(cfg *config.Config, memory memory.Memory) (*GoFetch, error) {
	return &GoFetch{
		cfg:        cfg,
		memory:     memory,
		downloader: collector.NewTransmissionDownloader(&cfg.Downloader),
	}, nil
}

func (gf *GoFetch) Fetch() (dls []Downloadable, err error) {
	for _, entry := range gf.cfg.Entries {
		source, ok := gf.cfg.Sources[entry.SourceName]
		if !ok {
			return nil, errors.New("source not found")
		}

		c, err := source.Collector()
		if err != nil {
			return nil, err
		}

		downloads, err := c.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		filteredDownloads, err := entry.Filter.Filter(downloads)
		if err != nil {
			return nil, err
		}

		for _, dl := range filteredDownloads {
			dls = append(dls, Downloadable{
				Downloadable: dl,
				Optional:     dl.Optional,
				Seen:         gf.memory.Has(collector.Hash(dl)),
			})
		}
	}

	return dls, nil
}

func (g *GoFetch) Download(dl Downloadable) error {
	hash := collector.Hash(dl)
	if g.memory.Has(hash) {
		return fmt.Errorf("already downloaded: %s", dl.Name())
	}

	err := g.downloader.Download(dl)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	err = g.memory.Put(hash, "d")
	if err != nil {
		return fmt.Errorf("failed to save download to memory: %w", err)
	}

	return nil
}

func (gf *GoFetch) Stream(dl Downloadable) error {
	downloader := collector.WebTorrentDownloader{}
	err := downloader.Download(dl)
	return err
}

func (gf *GoFetch) Ignore(dl Downloadable) error {
	hash := collector.Hash(dl)

	if gf.memory.Has(hash) {
		return fmt.Errorf("already ignored: %s", dl.Name())
	}

	err := gf.memory.Put(hash, "i")
	if err != nil {
		return fmt.Errorf("failed to save ignore to memory: %w", err)
	}

	return nil
}
