package collector

import (
	"context"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/roarc0/go-magnet"
)

type NyaaMagnetCollector struct {
	colly *colly.Collector
	uri   string
}

func NewNyaaMagnetCollector(uri string, opts ...CollectorOption) (*NyaaMagnetCollector, error) {
	c := &NyaaMagnetCollector{
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

func (c *NyaaMagnetCollector) Collect(ctx context.Context) ([]Downloadable, error) {
	dls := []Downloadable{}
	tmpMagnet := Magnet{}

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			title := e.Attr("title")
			if len(title) != 0 {
				tmpMagnet.name = title
				return // the title should appear before href
			}

			href := e.Attr("href")
			if strings.HasPrefix(href, "magnet:") {
				tmpMagnet.uri = href
			}

			magnet, err := magnet.Parse(href)
			if err != nil {
				return
			}
			if magnet.ExactLength != 0 {
				tmpMagnet.size = magnet.ExactLength
			}

			if tmpMagnet.name != "" && tmpMagnet.uri != "" {
				dls = append(dls, tmpMagnet)
				tmpMagnet = Magnet{}
			}
		})

	if err := c.colly.Visit(c.uri); err != nil {
		return nil, err
	}

	return dls, nil
}
