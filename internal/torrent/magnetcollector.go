package torrent

import (
	"context"
	"strings"

	"github.com/gocolly/colly"
	"github.com/roarc0/gct/internal/collector"
)

type MagnetCollector struct {
	colly *colly.Collector
	uri   string
}

func NewMagnetCollector(uri string) *MagnetCollector {
	return &MagnetCollector{
		uri:   uri,
		colly: colly.NewCollector(),
	}
}

func (c *MagnetCollector) Collect(ctx context.Context) ([]collector.Downloadable, error) {
	dls := []collector.Downloadable{}

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.HasPrefix(href, "magnet:") {
				dls = append(dls, Magnet{
					name: e.Text,
					uri:  e.Attr("href"),
				})
			}
		})

	err := c.colly.Visit(c.uri)
	if err != nil {
		return nil, err
	}

	return dls, nil
}
