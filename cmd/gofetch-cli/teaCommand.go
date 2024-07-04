package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/shlex"
	"github.com/rs/zerolog/log"

	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/gofetch"
)

type model struct {
	gf        *gofetch.GoFetch
	textInput textinput.Model
	err       error
}

func commandModel(gf *gofetch.GoFetch) model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.Focus()
	ti.ShowSuggestions = true
	ti.SetSuggestions([]string{"fetch", "help", "quit", "clear"})
	ti.CompletionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333"))
	ti.PromptStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))
	ti.CharLimit = 128
	ti.Width = 20

	return model{
		gf:        gf,
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyEnter:
			{
				args, err := shlex.Split(m.textInput.Value())
				if err != nil {
					return m, nil
				}

				log.Info().
					Any("args", args).
					Msg("Command")

				if len(args) == 0 {
					return m, nil
				}

				switch args[0] {
				case "fetch", "f":
					nm := newDownloadsModel(m.gf, func() ([]gofetch.Downloadable, error) {
						return m.gf.Fetch()
					})

					return nm, nm.Init()
				case "search", "s":
					if len(args) != 3 {
						return m, nil
					}

					nm := newDownloadsModel(m.gf, func() ([]gofetch.Downloadable, error) {
						filter := filter.NewFilter([]filter.Matcher{
							&filter.RegexMatcher{
								MatchType: filter.MatchTypeRequired,
								Regex:     args[2],
							},
						})
						return m.gf.Search(args[1], filter)
					})

					return nm, nm.Init()
				case "help", "h":
					//nm := newHelpModel(m.gf)
					//return nm, nm.Init()
				case "quit", "q":
					return m, tea.Quit
				case "clear", "c":
					return m, tea.ClearScreen
				}
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := ""
	if m.err != nil {
		s += m.err.Error() + "\n\n"
	}

	s += fmt.Sprintf(
		"GoFetch\n\n%s\n",
		m.textInput.View(),
	)

	return s
}
