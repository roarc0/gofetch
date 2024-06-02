package collector

import (
	"context"
	"strings"

	"github.com/gocolly/colly"
	"github.com/roarc0/go-magnet"
)

type MagnetCollector struct {
	colly *colly.Collector
	uri   string
}

func NewMagnetCollector(uri string) (*MagnetCollector, error) {
	return &MagnetCollector{
		uri:   uri,
		colly: colly.NewCollector(),
	}, nil
}

func (c *MagnetCollector) Collect(ctx context.Context) ([]Downloadable, error) {
	dls := []Downloadable{}

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.HasPrefix(href, "magnet:") {
				magnet, err := magnet.Parse(href)
				if err != nil {
					return
				}

				name := e.Text
				if len(magnet.DisplayNames) > 0 {
					name = magnet.DisplayNames[0]
				}

				dls = append(dls, Magnet{
					name: name,
					uri:  e.Attr("href"),
					size: magnet.ExactLength,
				})
			}
		})

	err := c.colly.Visit(c.uri)
	if err != nil {
		return nil, err
	}

	return dls, nil
}
