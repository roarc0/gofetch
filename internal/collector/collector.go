package collector

import (
	"context"
	"crypto"
	"fmt"
	"time"
)

// DownloadableCollector is an interface that defines a method to collect downloadables.
type DownloadableCollector interface {
	Collect(ctx context.Context) ([]Downloadable, error)
}

// Downloadable is an interface that defines a method to get the name and URI of a downloadable.
type Downloadable interface {
	fmt.Stringer
	Name() string
	URI() string
	Size() uint64
	Date() time.Time
}

// Downloader is an interface that defines a method to open a downloadable.
type Downloader interface {
	Download(d Downloadable) error
}

func Hash(d Downloadable) string {
	h := crypto.SHA1.New()
	h.Write([]byte(d.URI()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type CollectorConfig struct {
	HTTP      HttpConfig
	PageCount int
}

type CollectorOption func(*CollectorConfig) error

func WithPageCount(page int) CollectorOption {
	return func(cfg *CollectorConfig) error {
		cfg.PageCount = page
		return nil
	}
}

func WithHTTPConfig(http HttpConfig) CollectorOption {
	return func(cfg *CollectorConfig) error {
		cfg.HTTP = http
		return nil
	}
}

func processOptions(opts ...CollectorOption) (*CollectorConfig, error) {
	cfg := CollectorConfig{}
	for _, opt := range opts {
		err := opt(&cfg)
		if err != nil {
			return nil, err
		}
	}
	cfg.HTTP.SetDefaultsOnEmptyFields()
	return &cfg, nil
}
