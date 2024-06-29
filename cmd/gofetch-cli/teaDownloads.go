package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/roarc0/gofetch/internal/gofetch"
)

const (
	tableHeight = 30
)

type inDlsMsg []gofetch.Downloadable

type outDl struct {
	idx int
	err error
}

type outDlsMsg []outDl

type downloadsModel struct {
	gf *gofetch.GoFetch

	err error

	fetchingSpinner spinner.Model
	fetched         bool

	table       table.Model
	selections  map[int]gofetch.Action
	interactive bool

	dls []gofetch.Downloadable
}

func newDownloadsModel(gf *gofetch.GoFetch) tea.Model {
	m := downloadsModel{
		gf:         gf,
		selections: make(map[int]gofetch.Action, 0),
	}

	spn := spinner.New()
	spn.Spinner = spinner.Dot
	spn.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.fetchingSpinner = spn

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Action", Width: 10},
			{Title: "Name", Width: 80},
		}),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.table = t

	return m
}

func (m downloadsModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchingSpinner.Tick,
		m.fetchCommand(),
	)
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
			m.table.MoveUp(1)
		case "down", "j":
			m.table.MoveDown(1)
		case "ctrl+k":
			m.table.MoveUp(5)
		case "ctrl+j":
			m.table.MoveDown(5)
		case " ":
			m.selectionChange()
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
			// if !m.interactive || m.err != nil {
			// 	return commandModel(m.gf), nil
			// }
			return m, m.downloadCommand()
		}
	case inDlsMsg:
		m.fetchDone(msg)
	case errMsg:
		m.fetched = true
		m.err = msg
	case outDlsMsg:
		if len(m.dls) == 0 {
			return commandModel(m.gf), nil
		}
		nm := newDownloadsModel(m.gf)
		return nm, nm.Init()
	case spinner.TickMsg:
		if !m.fetched {
			var cmd tea.Cmd
			m.fetchingSpinner, cmd = m.fetchingSpinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m downloadsModel) View() string {
	if !m.fetched {
		return fmt.Sprintf("\n\n %s Loading ... press q to quit\n\n", m.fetchingSpinner.View())
	}

	if m.err != nil {
		return "\n" + m.err.Error() + "\n"
	}

	return m.showTable()
}

func (m downloadsModel) showTable() string {
	if len(m.dls) == 0 {
		return "\nNo items found.\n"
	}

	var s string
	s += "Press <j>,<k> or use arrow keys to navigate.\n"
	if m.interactive {
		s += "Press <space> to change the action to perform on the item \n"
		s += "Press <u> unselect all, <a> download all, <i> ignore all.\n"
	}
	s += "Press <Ctrl+s> to stream (needs webtorrent installed)\n\n"

	s += m.table.View() + "\n"

	s += "\nPress <enter> to proceed, <esc> to abort.\n"
	return s
}

func (m *downloadsModel) fetchDone(dls inDlsMsg) {
	m.fetched = true

	for _, dl := range dls {
		if dl.Action.Seen() {
			continue
		}

		i := len(m.dls)
		if dl.Optional {
			m.selections[i] = gofetch.IgnoreAction
		} else {
			m.selections[i] = gofetch.DownloadAction
		}

		m.dls = append(m.dls, dl)
	}

	m.interactive = len(m.dls) != 0

	if !m.interactive {
		m.dls = dls
		clear(m.selections)
	}

	m.updateTable()
}

func (m *downloadsModel) updateTable() {
	rows := []table.Row{}

	for i, dl := range m.dls {
		var actionStr string
		if !m.interactive {
			actionStr = fmt.Sprintf("[%s]", dl.Action.String())
		} else {
			actionStr = m.selections[i].String()
			if len(actionStr) != 0 {
				actionStr = fmt.Sprintf("<%s>", actionStr)
			}
		}

		rows = append(rows,
			table.Row{
				actionStr,
				dl.Name(),
			},
		)
	}

	m.table.SetRows(rows)
}

func (m downloadsModel) getAction(idx int) gofetch.Action {
	action := gofetch.NoAction
	if v, ok := m.selections[idx]; ok {
		action = v
	}
	return action
}

func (m downloadsModel) fetchCommand() tea.Cmd {
	return func() tea.Msg {
		dls, err := m.gf.Fetch()
		if err != nil {
			return errMsg(err)
		}

		return inDlsMsg(dls)
	}
}

func (m downloadsModel) downloadCommand() tea.Cmd {
	return func() tea.Msg {
		var msg outDlsMsg

		for i, action := range m.selections {
			var err error

			switch action {
			case gofetch.DownloadAction:
				err = m.gf.Download(m.dls[i])
			case gofetch.IgnoreAction:
				err = m.gf.Ignore(m.dls[i])
			}

			msg = append(msg, outDl{i, err})
		}

		return msg
	}
}

func (m downloadsModel) streamCommand() tea.Cmd {
	return func() tea.Msg {
		dl := m.dls[m.table.Cursor()]

		err := m.gf.Stream(dl)

		return errMsg(err)
	}
}

func (m *downloadsModel) selectionChangeAll(from gofetch.Action, to gofetch.Action) {
	if !m.interactive {
		return
	}

	for i := range m.dls {
		if from == gofetch.NoAction {
			m.selections[i] = to
			continue
		}
		if to == gofetch.NoAction {
			delete(m.selections, i)
			continue
		}
		if v, ok := m.selections[i]; ok && v == from {
			m.selections[i] = to
		}
	}

	m.updateTable()
}

func (m *downloadsModel) selectionChange() {
	if !m.interactive {
		return
	}

	cursor := m.table.Cursor()
	if v, ok := m.selections[cursor]; ok {
		switch v {
		case gofetch.DownloadAction:
			m.selections[cursor] = gofetch.IgnoreAction
		case gofetch.IgnoreAction:
			delete(m.selections, cursor)
		}
	} else {
		m.selections[cursor] = gofetch.DownloadAction
	}

	m.updateTable()
}
