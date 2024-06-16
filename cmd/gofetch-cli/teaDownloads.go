package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

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

type downloadsModel struct {
	gf *gofetch.GoFetch

	allDls []gofetch.Downloadable
	newDls []gofetch.Downloadable

	err      error
	cursor   int
	selected map[int]action

	fetched bool
}

type dlsMsg []gofetch.Downloadable

type dlDoneMsg struct {
	index int
	error error
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
		if len(m.newDls) == 0 {
			return commandModel(m.gf), nil
		}
		nm := newDownloadsModel(m.gf)
		return nm, nm.Init()
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
	s := "Press <space> to change the action to perform on each item \n"
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
		s += showListCapped(m.allDls)
	}
	s += "\n"
	return s
}

func showListCapped(dls []gofetch.Downloadable) string {
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
	m.allDls = dls

	for _, dl := range dls {
		if dl.Seen {
			continue
		}
		m.newDls = append(m.newDls, dl)
	}

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

		return dlsMsg(dls)
	}
}

func (m downloadsModel) downloadCommand() tea.Cmd {
	return func() tea.Msg {
		var doneMsg []dlDoneMsg
		for i, action := range m.selected {
			dl := m.newDls[i]
			var err error
			switch action {
			case downloadAction:
				err = m.gf.Download(dl)
			case ignoreAction:
				err = m.gf.Ignore(dl)
			}

			doneMsg = append(doneMsg, dlDoneMsg{i, err})
		}

		return doneMsg
	}
}
