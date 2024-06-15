package collector

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-retryablehttp"
)

var (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
)

func newColly() *colly.Collector {
	c := colly.NewCollector()

	c.SetClient(httpClient())
	c.UserAgent = userAgent

	return c
}

func httpClient() *http.Client {
	var (
		dnsResolverIP        = "1.1.1.1:53"
		dnsResolverProto     = "udp"
		dnsResolverTimeoutMs = 10000
	)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
	}
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	tr := &http.Transport{
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: dialContext,
	}
	retryClient.HTTPClient.Transport = tr
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(l retryablehttp.Logger, req *http.Request, attempt int) {
		// remove Accept header since it can cause issues with some websites
		req.Header.Del("Accept")
	}

	return retryClient.StandardClient()
}
