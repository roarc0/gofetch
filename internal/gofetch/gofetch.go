package gofetch

import (
	"context"
	"errors"
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
	cfg    *config.Config
	memory memory.Memory
}

// NewGoFetch creates a new GoFetch object.
func NewGoFetch(cfg *config.Config, memory memory.Memory) (*GoFetch, error) {
	return &GoFetch{
		cfg:    cfg,
		memory: memory,
	}, nil
}

func (g *GoFetch) Fetch() (dls []filter.MatchedDownloadable, err error) {
	for _, entry := range g.cfg.Entries {
		source, ok := g.cfg.Sources[entry.SourceName]
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

func (g *GoFetch) Download(dls []filter.MatchedDownloadable, filterOptional bool) {
	var filteredDls []collector.Downloadable
	if filterOptional {
		filteredDls = filter.FilterDownloadables(dls, nil)
	} else {
		filteredDls = filter.FilterDownloadables(dls, func(filter.MatchedDownloadable) bool { return false })
	}

	downloader := collector.NewTransmissionDownloader(nil, &g.cfg.Transmission)

	for _, dl := range filteredDls {
		hash := collector.Hash(dl)

		if g.memory.Has(hash) {
			log.Printf("Already downloaded %q\n", dl.Name())
			continue
		}

		log.Printf("Downloading: %q\n", dl.Name())

		downloader.Downloadable = dl
		err := downloader.Download()
		if err != nil {
			log.Println(err)
		}

		err = g.memory.Put(hash)
		if err != nil {
			log.Println(err)
		}
	}
}
