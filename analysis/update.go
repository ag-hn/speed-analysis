package analysis

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ag-hn/speed-analysis/polish"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.Disabled {
		return m, nil
	}

	switch msg := msg.(type) {
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	case errorMsg:
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(polish.Colors.Red600).
				Bold(true).
				Render(string(msg))))
	case copyToClipboardMsg:
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Bold(true).
				Render(string(msg))))
	case statusMessageTimeoutMsg:
		m.StatusMessage = ""
	case waitForProcessedItemMsg:
		m.count++
		m.processed = append(m.processed, msg.processed)
		return m, m.waitForProcessedItem()
	case listenForProcessedItemMsg:
		m.State = IdleState
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Down):
			if m.State != IdleState {
				return m, nil
			}

			m.Cursor++

			if m.Cursor >= len(m.processed) {
				m.Cursor = len(m.processed) - 1
			}

			if m.Cursor > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.keyMap.Up):
			if m.State != IdleState {
				return m, nil
			}

			m.Cursor--

			if m.Cursor < 0 {
				m.Cursor = 0
			}

			if m.Cursor < m.min {
				m.min--
				m.max--
			}	
		case key.Matches(msg, m.keyMap.CopyPathToClipboard):
			if m.State != IdleState {
				return m, nil
			}

			return m, copyToClipboardCmd(m.processed[m.Cursor])
		// case key.Matches(msg, m.keyMap.Process):
		// 	m.State = ProcessingState
		// 	cmds = append(cmds, m.ListenForProcessedItem()) // generate activity
		}
	}

	return m, tea.Batch(cmds...)
}
