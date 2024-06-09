package main

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/gofetch"
)

type downloadsModel struct {
	gf *gofetch.GoFetch

	newDls []filter.MatchedDownloadable
	allDls []filter.MatchedDownloadable

	err      error
	cursor   int
	selected map[int]struct{}

	fetched bool
}

type dlsMsg struct {
	new []filter.MatchedDownloadable
	all []filter.MatchedDownloadable
}

type dlsDoneMsg []collector.Downloadable

func (m downloadsModel) Init() tea.Cmd {
	return fetch(m.gf)
}

func (m downloadsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
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
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			if m.err != nil {
				return commandModel(m.gf), nil
			}
			return m, m.download()
		}
	case dlsMsg:
		m.fetchDone(msg)
	case errMsg:
		m.fetched = true
		m.err = msg
	case dlsDoneMsg:
		return commandModel(m.gf), nil
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
		s := "\nNo new downloads found.\n"
		if len(m.allDls) > 0 {
			s += fmt.Sprintf("Already downloaded count: %d\n", len(m.allDls))
			max := 10
			if len(m.allDls) < max {
				max = len(m.allDls)
			}

			for _, dl := range m.allDls[:max] {
				s += fmt.Sprintf("- %s\n", dl.Name())
			}
			s += "...\n"
		}
		s += "\n"
		return s
	}

	s := "Select the items you want to download \n"

	for i, dl := range m.newDls {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, dl.Name())
	}

	s += "\nPress <enter> to download, <esc> to abort.\n"

	return s
}

func newDownloadsModel(gf *gofetch.GoFetch) tea.Model {
	return downloadsModel{
		gf:       gf,
		selected: make(map[int]struct{}, 0),
	}
}

func fetch(gf *gofetch.GoFetch) tea.Cmd {
	return func() tea.Msg {
		dls, err := gf.Fetch()
		if err != nil {
			return errMsg(err)
		}

		undownloadedDls := gf.Undownloaded(dls)

		return dlsMsg{
			new: undownloadedDls,
			all: dls,
		}
	}
}

func (m *downloadsModel) fetchDone(dls dlsMsg) {
	m.fetched = true
	m.newDls = dls.new
	m.allDls = dls.all
	for i, dl := range m.newDls {
		if !dl.Optional {
			m.selected[i] = struct{}{}
		}
	}
}

func (m downloadsModel) download() tea.Cmd {
	return func() tea.Msg {
		dls := make([]collector.Downloadable, 0)
		for i := range m.selected {
			dls = append(dls, m.newDls[i].Downloadable)
		}
		var errs []error
		for _, dl := range dls {
			if err := m.gf.Download(dl); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return errMsg(errors.Join(errs...))
		}
		return dlsDoneMsg(dls)
	}
}