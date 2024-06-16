package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"

	"github.com/roarc0/gofetch/internal/gofetch"
)

type (
	errMsg error
)

func runTea(gf *gofetch.GoFetch) {
	p := tea.NewProgram(commandModel(gf))
	if _, err := p.Run(); err != nil {
		log.Error().Err(err).Msg("Alas, there's been an error")
		os.Exit(1)
	}
}
