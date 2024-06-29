package collector

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultUserAgent            = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	defaultDNSResolverIP        = "1.1.1.1:53"
	defaultDNSResolverProto     = "udp"
	defaultDNSResolverTimeoutMs = 5000
	defaultRetryMax             = 2
)

type HttpConfig struct {
	UserAgent   string
	DNSResolver string
	DNSProto    string
	DNSTimeout  int
	Insecure    bool
	RetryMax    int
}

func (cfg *HttpConfig) SetDefaultsOnEmptyFields() {
	if cfg.DNSResolver == "" {
		cfg.DNSResolver = defaultDNSResolverIP
	}
	if cfg.DNSProto == "" {
		cfg.DNSProto = defaultDNSResolverProto
	}
	if cfg.DNSTimeout == 0 {
		cfg.DNSTimeout = defaultDNSResolverTimeoutMs
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = defaultUserAgent
	}
	if cfg.RetryMax == 0 {
		cfg.RetryMax = defaultRetryMax
	}
}

func newColly(cfg *HttpConfig) *colly.Collector {
	cfg.SetDefaultsOnEmptyFields()

	c := colly.NewCollector()
	c.UserAgent = defaultUserAgent
	c.SetClient(httpClient(cfg))

	return c
}

func httpClient(cfg *HttpConfig) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = cfg.RetryMax
	retryClient.HTTPClient.Timeout = 7 * time.Second
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(cfg.DNSTimeout) * time.Millisecond,
				}
				return d.DialContext(ctx, cfg.DNSProto, cfg.DNSResolver)
			},
		},
	}
	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	tr := &http.Transport{
		DialContext: dialContext,
	}
	if cfg.Insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	retryClient.HTTPClient.Transport = tr
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(l retryablehttp.Logger, req *http.Request, attempt int) {
		req.Header.Del("Accept")
	}

	return retryClient.StandardClient()
}
