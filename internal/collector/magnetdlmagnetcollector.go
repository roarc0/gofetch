package collector

import (
	"context"
	"strings"

	"github.com/gocolly/colly/v2"
)

type MagnetDLMagnetCollector struct {
	colly *colly.Collector
	uri   string
}

func NewMagnetDLMagnetCollector(uri string, opts ...CollectorOption) (*MagnetDLMagnetCollector, error) {
	c := &MagnetDLMagnetCollector{
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

func (c *MagnetDLMagnetCollector) Collect(ctx context.Context) ([]Downloadable, error) {
	dls := []Downloadable{}
	tmpMagnet := Magnet{}

	c.colly.OnHTML("a",
		func(e *colly.HTMLElement) {
			href := e.Attr("href")
			if strings.HasPrefix(href, "magnet:") {
				tmpMagnet.uri = href
				return
			}

			if len(e.Attr("title")) != 0 {
				tmpMagnet.name = e.Attr("title")
			}

			if tmpMagnet.name != "" && tmpMagnet.uri != "" {
				dls = append(dls, tmpMagnet)
				tmpMagnet = Magnet{}
			}
		})

	// c.colly.OnHTML("td",
	// 	func(e *colly.HTMLElement) {
	// 		text := e.Text

	// 		if strings.HasSuffix(text, "MB") {
	// 			val, err := strconv.Atoi(strings.Split(text, " ")[0])
	// 			if err == nil {
	// 				tmpMagnet.size = uint64(val)
	// 			}
	// 		}

	// 		if strings.HasSuffix(text, "months") {
	// 			t, err := dateparse.ParseAny(fmt.Sprintf("%s ago", text))
	// 			if err != nil {
	// 				return
	// 			}
	// 			tmpMagnet.time = t
	// 		}

	// 		if tmpMagnet.name != "" && tmpMagnet.uri != "" {
	// 			dls = append(dls, tmpMagnet)
	// 			tmpMagnet = Magnet{}
	// 		}
	// 	})

	if err := c.colly.Visit(c.uri); err != nil {
		return nil, err
	}

	return dls, nil
}
