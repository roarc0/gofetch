package gct

import (
	"context"
	"log"

	"github.com/roarc0/gct/internal/collector"
	"github.com/roarc0/gct/internal/filter"
	"github.com/roarc0/gct/internal/torrent"
)

type GCT struct {
	Sources map[string]Source `json:"sources"`
	Entries map[string]Entry  `json:"entries"`
}

type Entry struct {
	Source string        `json:"source"`
	Filter filter.Filter `json:"filter"`
}

type Source struct {
	Name string   `json:"name"`
	URIs []string `json:"uris"`
}

// NewGCT creates a new GCT object.
func NewGCT() *GCT {
	return &GCT{
		Sources: map[string]Source{
			"nyaa": {
				Name: "nyaa",
				URIs: []string{"https://nyaa.si/?c=1_2&s=seeders&o=desc"},
			},
		},
		Entries: map[string]Entry{
			"Kaijuu 8": {
				Source: "nyaa",
				Filter: filter.NewFilterWithOptionalMatches([]filter.Matcher{
					&filter.RegexMatcher{Regex: ".*Kaijuu 8.*"},
				}),
			},
		},
	}
}

func (g *GCT) Fetch() error {
	for _, entry := range g.Entries {
		var collector collector.DownloadableCollector
		if entry.Source == "nyaa" {
			uris := g.Sources["nyaa"].URIs
			c, err := torrent.NewNyaaMagnetCollector(uris[0])
			if err != nil {
				return err
			}
			collector = c
		}

		dls, err := collector.Collect(context.Background())
		if err != nil {
			return err
		}

		for _, e := range g.Entries {
			if e.Filter != nil {
				matched, err := e.Filter.Filter(dls)
				if err != nil {
					return err
				}
				dls = filter.FilterPartialMatchDownloadables(matched)
			}
		}

		log.Printf("ToDownload:\n%v\n", dls)
	}

	return nil
}
