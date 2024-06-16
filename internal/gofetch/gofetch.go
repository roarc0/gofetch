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
	Action   Action
}

type Action int

const (
	NoAction Action = iota
	DownloadAction
	IgnoreAction
)

func (a Action) String() string {
	switch a {
	case NoAction:
		return ""
	case DownloadAction:
		return "downloaded"
	case IgnoreAction:
		return "ignored"
	default:
		return "unknown"
	}
}

func (a Action) Seen() bool {
	return a == DownloadAction || a == IgnoreAction
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
			var action Action

			actionPtr, err := gf.memory.Get(collector.Hash(dl))
			if err != nil {
				return nil, err
			}
			switch *actionPtr {
			case DownloadAction.String():
				action = DownloadAction
			case IgnoreAction.String():
				action = IgnoreAction
			case NoAction.String():
				fallthrough
			default:
				action = NoAction
			}

			dls = append(dls, Downloadable{
				Downloadable: dl,
				Optional:     dl.Optional,
				Action:       action,
			})
		}
	}

	return dls, nil
}

func (g *GoFetch) Download(dl Downloadable) error {
	hash := collector.Hash(dl)
	if g.memory.Has(hash) {
		return fmt.Errorf("already processed: %s", dl.Name())
	}

	err := g.downloader.Download(dl)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	err = g.memory.Put(hash, DownloadAction.String())
	if err != nil {
		return fmt.Errorf("failed to save to memory: %w", err)
	}

	return nil
}

func (gf *GoFetch) Ignore(dl Downloadable) error {
	hash := collector.Hash(dl)

	if gf.memory.Has(hash) {
		return fmt.Errorf("already ignored: %s", dl.Name())
	}

	err := gf.memory.Put(hash, IgnoreAction.String())
	if err != nil {
		return fmt.Errorf("failed to save ignore to memory: %w", err)
	}

	return nil
}

func (gf *GoFetch) Stream(dl Downloadable) error {
	return collector.WebTorrentDownloader{}.Download(dl)
}
