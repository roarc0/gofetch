package main

import (
	"flag"

	"github.com/roarc0/gofetch/internal/config"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/gofetch"
	"github.com/roarc0/gofetch/internal/logger"
	"github.com/roarc0/gofetch/internal/memory"

	"github.com/rs/zerolog/log"
)

func main() {
	cfgPath := flag.String("config", ".", "Path to the configuration file")
	flag.Parse()

	logger.SetupLogger()

	cfg, err := config.LoadYaml(*cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("config.LoadYaml()")
	}

	memory, err := memory.NewMemory(cfg.Memory.FilePath, "downloads")
	if err != nil {
		log.Fatal().Err(err).Msg("memory.NewMemory()")
	}

	gf, err := gofetch.NewGoFetch(cfg, memory)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GoFetch object")
	}

	dls, err := gf.Fetch()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch")
	}

	gf.DownloadAll(filter.FilterOptionalMatches(dls))

	if err != nil {
		log.Error().Err(err).Msg("Failed to create GoFetch object")
	}
}
