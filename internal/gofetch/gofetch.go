package gofetch

import (
	"context"
	"crypto"
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

func (g *GoFetch) DownloadAll(dls []collector.Downloadable) {
	for _, dl := range dls {
		h := crypto.SHA1.New()
		h.Write([]byte(dl.URI()))
		hex := h.Sum(nil)

		if g.memory.Has(string(hex)) {
			log.Printf("Already downloaded %q\n", dl.Name())
			continue
		}

		log.Printf("Downloading %q\n", dl.Name())

		err := collector.XDGDownloader{
			Downloadable: dl,
		}.Download()
		if err != nil {
			log.Println(err)
		}

		err = g.memory.Put(string(hex))
		if err != nil {
			log.Println(err)
		}
	}
}
