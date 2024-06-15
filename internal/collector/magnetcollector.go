package collector

import (
	"context"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/roarc0/go-magnet"
)

type MagnetCollector struct {
	colly *colly.Collector
	uri   string
}

func NewMagnetCollector(uri string, opts ...CollectorOption) (*MagnetCollector, error) {
	c := &MagnetCollector{
		uri: uri,
	}

	cfg, err := processOptions(opts...)
	if err != nil {
		return nil, err
	}

	if c.colly == nil {
		c.colly = newColly(&cfg.HTTP)
	}

	return c, nil
}

func (c *MagnetCollector) Collect(ctx context.Context) ([]Downloadable, error) {
	dls := []Downloadable{}

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if !strings.HasPrefix(href, "magnet:?") {
				return
			}

			magnet, err := magnet.Parse(href)
			if err != nil {
				return
			}

			name := e.Text
			if len(magnet.DisplayNames) > 0 {
				name = magnet.DisplayNames[0]
			}

			dls = append(dls,
				Magnet{
					name: name,
					uri:  href,
					size: magnet.ExactLength,
				})
		})

	if err := c.colly.Visit(c.uri); err != nil {
		return nil, err
	}

	return dls, nil
}
