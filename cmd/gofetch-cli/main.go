package main

import (
	"flag"

	"github.com/roarc0/gofetch/internal/config"
	"github.com/roarc0/gofetch/internal/gofetch"
	"github.com/roarc0/gofetch/internal/logger"
	"github.com/roarc0/gofetch/internal/memory"

	"github.com/rs/zerolog/log"
)

func main() {
	cfgPath := flag.String("config", ".", "Path to the configuration file")
	mode := flag.String("mode", "manual", "Mode to run the application in")
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

	switch *mode {
	case "auto":
		autoFetchAll(gf, true)
	case "manual":
		runTea(gf)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to create GoFetch object")
	}
}
