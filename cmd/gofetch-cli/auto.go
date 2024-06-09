package main

import (
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/gofetch"
	"github.com/rs/zerolog/log"
)

func fetchAll(gf *gofetch.GoFetch, filterOptional bool) {
	dls, err := gf.Fetch()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch")
	}

	var filteredDls []filter.MatchedDownloadable
	if filterOptional {
		filteredDls = filter.FilterDownloadables(dls, nil)
	} else {
		filteredDls = filter.FilterDownloadables(dls, func(filter.MatchedDownloadable) bool { return false })
	}

	gf.DownloadAll(filter.ToDownloadables(filteredDls))
}
