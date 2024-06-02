package torrent

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
	"github.com/roarc0/gct/internal/collector"
)

type NyaaMagnetCollector struct {
	colly *colly.Collector
	uri   string
}

type NyaaMagnetCollectorOption func(*NyaaMagnetCollector) error

func WithPage(page int) NyaaMagnetCollectorOption {
	return func(c *NyaaMagnetCollector) error {
		url, err := url.Parse(c.uri)
		if err != nil {
			return errors.Wrap(err, "failed to parse url in WithPage")
		}

		q := url.Query()
		q.Set("p", fmt.Sprint(page))
		url.RawQuery = q.Encode()

		c.uri = url.String()

		return nil
	}
}

func WithColly(colly *colly.Collector) NyaaMagnetCollectorOption {
	return func(c *NyaaMagnetCollector) error {
		c.colly = colly
		return nil
	}
}

func NewNyaaMagnetCollector(uri string, opts ...NyaaMagnetCollectorOption) (*NyaaMagnetCollector, error) {
	c := &NyaaMagnetCollector{
		uri: uri,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.colly == nil {
		colly.NewCollector()
	}

	return c, nil
}

func (c *NyaaMagnetCollector) Collect(ctx context.Context) ([]collector.Downloadable, error) {
	dls := []collector.Downloadable{}
	tmpName := ""

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			if len(e.Attr("title")) != 0 {
				tmpName = e.Attr("title")
				return
			}

			href := e.Attr("href")
			if strings.HasPrefix(href, "magnet:") {
				dls = append(dls, Magnet{
					name: tmpName,
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
