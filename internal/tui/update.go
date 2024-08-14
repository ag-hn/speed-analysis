package tui

import (
	"github.com/ag-hn/speed-analysis/analysis"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		halfSize := msg.Width / 2
		height := msg.Height

		m.analysis.SetSize(halfSize, height-3)
		m.help.SetSize(halfSize, height)

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Quit):
			if m.analysis.State != analysis.ProcessingState {
				return m, tea.Quit
			}
		}
	}

	m.analysis, cmd = m.analysis.Update(msg)
	cmds = append(cmds, cmd)
	
	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
