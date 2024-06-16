package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/gofetch"
)

const (
	maxPrintableDownloads = 100
)

type action int

const (
	downloadAction action = iota
	ignoreAction
)

func (a action) String() string {
	switch a {
	case downloadAction:
		return "download"
	case ignoreAction:
		return "ignore"
	default:
		return ""
	}
}

type downloadsModel struct {
	gf *gofetch.GoFetch

	newDls []filter.MatchedDownloadable
	allDls []filter.MatchedDownloadable

	err      error
	cursor   int
	selected map[int]action

	fetched bool
}

type dlsMsg struct {
	new []filter.MatchedDownloadable
	all []filter.MatchedDownloadable
}

type dlDoneMsg struct {
	dl     collector.Downloadable
	action action
	error  error
}

func (m downloadsModel) Init() tea.Cmd {
	return fetchCommand(m.gf)
}

func (m downloadsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit
		case "esc":
			return commandModel(m.gf), nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.newDls)-1 {
				m.cursor++
			}
		case " ":
			m.updateSelected()
		case "enter":
			if m.err != nil {
				return commandModel(m.gf), nil
			}
			return m, m.downloadCommand()
		}
	case dlsMsg:
		m.fetchDone(msg)
	case errMsg:
		m.fetched = true
		m.err = msg
	case []dlDoneMsg:
		return commandModel(m.gf), nil
	}

	return m, nil
}

func (m *downloadsModel) updateSelected() {
	if v, ok := m.selected[m.cursor]; ok {
		switch v {
		case downloadAction:
			m.selected[m.cursor] = ignoreAction
		case ignoreAction:
			delete(m.selected, m.cursor)
		}
	} else {
		m.selected[m.cursor] = downloadAction
	}
}

func (m downloadsModel) View() string {
	if !m.fetched {
		return "\nFetching...\n"
	}

	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}

	if len(m.newDls) == 0 {
		return m.alreadyDownloadedList()
	}

	return m.downloadPrompt()
}

func (m downloadsModel) downloadPrompt() string {
	s := "Select the items you want to download or ignore \n"
	s += "[D] to download, [I] to ignore [ ] to do nothing\n\n"
	for i, dl := range m.newDls {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if v, ok := m.selected[i]; ok {
			switch v {
			case downloadAction:
				checked = "D"
			case ignoreAction:
				checked = "I"
			}
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, dl.Name())
	}

	s += "\nPress <enter> to proceed, <esc> to abort.\n"
	return s
}

func (m downloadsModel) alreadyDownloadedList() string {
	s := "\nNo new downloads found.\n"
	if len(m.allDls) > 0 {
		s += fmt.Sprintf("Already downloaded items count: %d\n", len(m.allDls))
		s += downloadsListCapped(m.allDls)
	}
	s += "\n"
	return s
}

func downloadsListCapped(dls []filter.MatchedDownloadable) string {
	max := maxPrintableDownloads
	if len(dls) < max {
		max = len(dls)
	}

	s := ""
	for _, dl := range dls[:max] {
		s += fmt.Sprintf("- %s\n", dl.Name())
	}
	if len(dls) > max {
		s += fmt.Sprintf("%d more ...\n", len(dls)-max)
	}
	return s
}

func newDownloadsModel(gf *gofetch.GoFetch) tea.Model {
	return downloadsModel{
		gf:       gf,
		selected: make(map[int]action, 0),
	}
}

func (m *downloadsModel) fetchDone(dls dlsMsg) {
	m.fetched = true
	m.newDls = dls.new
	m.allDls = dls.all
	for i, dl := range m.newDls {
		if !dl.Optional {
			m.selected[i] = downloadAction
		}
	}
}

func fetchCommand(gf *gofetch.GoFetch) tea.Cmd {
	return func() tea.Msg {
		dls, err := gf.Fetch()
		if err != nil {
			return errMsg(err)
		}

		undownloadedDls := gf.FilterNewDonwloads(dls)

		return dlsMsg{
			new: undownloadedDls,
			all: dls,
		}
	}
}

func (m downloadsModel) downloadCommand() tea.Cmd {
	return func() tea.Msg {
		var done []dlDoneMsg
		for _, i := range m.selected {
			dl := m.newDls[i].Downloadable
			var err error
			switch i {
			case downloadAction:
				err = m.gf.Download(dl)
			case ignoreAction:
				err = m.gf.Ignore(dl)
			}

			done = append(done, dlDoneMsg{dl, i, err})
		}

		return done
	}
}
