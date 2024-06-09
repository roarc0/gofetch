package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roarc0/gofetch/internal/gofetch"
)

type (
	errMsg error
)

func runTea(gf *gofetch.GoFetch) {
	p := tea.NewProgram(commandModel(gf))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
