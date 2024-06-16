package main

import (
	"github.com/rs/zerolog/log"

	"github.com/roarc0/gofetch/internal/gofetch"
)

func downloadAllNew(gf *gofetch.GoFetch) {
	dls, err := gf.Fetch()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch")
	}

	log.Info().Msgf("Found %d downloadables", len(dls))

	for _, dl := range dls {
		if dl.Optional {
			continue
		}
		err := gf.Download(dl)
		if err != nil {
			log.Error().Err(err).Msg("Error")
			continue
		}
		log.Info().Str("name", dl.Name()).Msg("Downloading")
	}
}
