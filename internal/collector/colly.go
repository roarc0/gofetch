package collector

import (
	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-retryablehttp"
)

func newColly() *colly.Collector {
	c := colly.NewCollector()

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	c.SetClient(retryClient.StandardClient())

	return c
}
