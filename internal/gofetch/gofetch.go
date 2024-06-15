package gofetch

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/config"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/memory"
)

const (
	downloadsBucket = "downloads"
)

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
		downloader: collector.NewTransmissionDownloader(&cfg.Transmission),
	}, nil
}

func (gf *GoFetch) Fetch() (dls []filter.MatchedDownloadable, err error) {
	for _, entry := range gf.cfg.Entries {
		source, ok := gf.cfg.Sources[entry.SourceName]
		if !ok {
			return nil, errors.New("source not found")
		}

		collector, err := source.Collector()
		if err != nil {
			return nil, err
		}

		d, err := collector.Collect(context.Background())
		if err != nil {
			return nil, err
		}

		fd, err := entry.Filter.Filter(d)
		if err != nil {
			return nil, err
		}

		dls = append(dls, fd...)
	}

	return dls, nil
}

func (gf *GoFetch) FilterNewDonwloads(dls []filter.MatchedDownloadable) []filter.MatchedDownloadable {
	return filter.FilterDownloadables(
		dls,
		func(d filter.MatchedDownloadable) bool {
			return gf.memory.Has(collector.Hash(d))
		},
	)
}

func (gf *GoFetch) DownloadAll(dls []collector.Downloadable) {
	for _, dl := range dls {
		err := gf.Download(dl)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		log.Println("Downloading:", dl.Name())
	}
}

func (g *GoFetch) Download(dl collector.Downloadable) error {
	hash := collector.Hash(dl)

	if g.memory.Has(hash) {
		return fmt.Errorf("already downloaded: %s", dl.Name())
	}

	err := g.downloader.Download(dl)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	err = g.memory.Put(hash)
	if err != nil {
		return fmt.Errorf("failed to save download to memory: %w", err)
	}

	return nil
}
