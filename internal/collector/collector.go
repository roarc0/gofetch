package collector

import (
	"context"
)

// DownloadableCollector is an interface that defines a method to collect downloadables.
type DownloadableCollector interface {
	Collect(ctx context.Context) ([]Downloadable, error)
}

// Downloadable is an interface that defines a method to get the name and URI of a downloadable.
type Downloadable interface {
	Name() string
	URI() string
	String() string
}

// Downloader is an interface that defines a method to open a downloadable.
type Downloader interface {
	Downloadable
	Download() error
}
