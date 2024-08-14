package analysis

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.ListenForProcessedItem(), // generate activity
		m.waitForProcessedItem(),   // wait for activity
	)	
}
