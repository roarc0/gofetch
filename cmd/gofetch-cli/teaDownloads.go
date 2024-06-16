package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/roarc0/gofetch/internal/gofetch"
)

const (
	maxPrintableDownloads = 100
)

type downloadsModel struct {
	gf *gofetch.GoFetch

	err      error
	cursor   int
	selected map[int]gofetch.Action

	fetched bool
	allDls  []gofetch.Downloadable
	newDls  []gofetch.Downloadable
}

type dlsMsg []gofetch.Downloadable

type dlDoneMsg struct {
	idx int
	err error
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
		case "esc", "q":
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
			m.selectionUpdate()
		case "a":
			m.selectionChangeAll(gofetch.NoAction, gofetch.DownloadAction)
			m.selectionChangeAll(gofetch.IgnoreAction, gofetch.DownloadAction)
		case "i":
			m.selectionChangeAll(gofetch.DownloadAction, gofetch.IgnoreAction)
			m.selectionChangeAll(gofetch.NoAction, gofetch.IgnoreAction)
		case "u":
			m.selectionChangeAll(gofetch.DownloadAction, gofetch.NoAction)
		case "ctrl+s":
			return m, m.streamCommand()
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

func (m downloadsModel) View() string {
	if !m.fetched {
		return "\nFetching...\n"
	}

	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}

	if len(m.newDls) == 0 {
		return m.showAlreadyProcessed()
	}

	return m.promptDownload()
}

func (m downloadsModel) promptDownload() string {
	s := "Press <space> to change the action to perform on each item \n"
	s += "[D] => download, [I] => ignore [ ] => do nothing\n\n"
	s += "press <j>,<k> or use arrow keys to navigate.\n"
	s += "press <u> unselect all, <a> download all, <i> ignore all.\n"
	s += "press <Ctrl+s> to stream (needs webtorrent installed)\n\n"

	for i, dl := range m.newDls {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if v, ok := m.selected[i]; ok {
			switch v {
			case gofetch.DownloadAction:
				checked = "D"
			case gofetch.IgnoreAction:
				checked = "I"
			}
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, dl.Name())
	}

	s += "\nPress <enter> to proceed, <esc> to abort.\n"
	return s
}

func (m downloadsModel) showAlreadyProcessed() string {
	if len(m.allDls) == 0 {
		return "\nNo items found.\n"
	}

	s := fmt.Sprintf("\nAll found items (%d) have been processed.\n", len(m.allDls))
	if len(m.allDls) > 0 {
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
		s += fmt.Sprintf("- [%s] %s\n", dl.Action, dl.Name())
	}
	if len(dls) > max {
		s += fmt.Sprintf("%d more ...\n", len(dls)-max)
	}
	return s
}

func newDownloadsModel(gf *gofetch.GoFetch) tea.Model {
	return downloadsModel{
		gf:       gf,
		selected: make(map[int]gofetch.Action, 0),
	}
}

func (m *downloadsModel) fetchDone(dls dlsMsg) {
	m.fetched = true
	m.allDls = dls

	for _, dl := range dls {
		if dl.Action.Seen() {
			continue
		}
		m.newDls = append(m.newDls, dl)
	}

	for i, dl := range m.newDls {
		if !dl.Optional {
			m.selected[i] = gofetch.DownloadAction
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
			case gofetch.DownloadAction:
				err = m.gf.Download(dl)
			case gofetch.IgnoreAction:
				err = m.gf.Ignore(dl)
			}

			doneMsg = append(doneMsg, dlDoneMsg{i, err})
		}

		return doneMsg
	}
}

func (m downloadsModel) streamCommand() tea.Cmd {
	return func() tea.Msg {
		dl := m.newDls[m.cursor]

		err := m.gf.Stream(dl)

		return errMsg(err)
	}
}

func (m *downloadsModel) selectionChangeAll(from gofetch.Action, to gofetch.Action) {
	for i := range m.newDls {
		if from == gofetch.NoAction {
			m.selected[i] = to
			continue
		}
		if to == gofetch.NoAction {
			delete(m.selected, i)
			continue
		}
		if v, ok := m.selected[i]; ok && v == from {
			m.selected[i] = to
		}
	}
}

func (m *downloadsModel) selectionUpdate() {
	if v, ok := m.selected[m.cursor]; ok {
		switch v {
		case gofetch.DownloadAction:
			m.selected[m.cursor] = gofetch.NoAction
		case gofetch.IgnoreAction:
			delete(m.selected, m.cursor)
		}
	} else {
		m.selected[m.cursor] = gofetch.DownloadAction
	}
}
