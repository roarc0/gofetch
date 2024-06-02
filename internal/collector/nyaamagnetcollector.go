package collector

import (
	"context"
	"strings"

	"github.com/gocolly/colly"
	"github.com/roarc0/go-magnet"
)

type NyaaMagnetCollector struct {
	colly *colly.Collector
	uri   string
}

type NyaaMagnetCollectorOption func(*NyaaMagnetCollector) error

// func WithPage(page int) NyaaMagnetCollectorOption {
// 	return func(c *NyaaMagnetCollector) error {
// 		url, err := url.Parse(c.uri)
// 		if err != nil {
// 			return errors.Wrap(err, "failed to parse url in WithPage")
// 		}

// 		q := url.Query()
// 		q.Set("p", fmt.Sprint(page))
// 		url.RawQuery = q.Encode()

// 		c.uri = url.String()

// 		return nil
// 	}
// }

// func WithColly(colly *colly.Collector) NyaaMagnetCollectorOption {
// 	return func(c *NyaaMagnetCollector) error {
// 		c.colly = colly
// 		return nil
// 	}
// }

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
		c.colly = colly.NewCollector()
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

	err := c.colly.Visit(c.uri)
	if err != nil {
		return nil, err
	}

	return dls, nil
}
