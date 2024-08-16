package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (m model) View() string {
	leftBox := m.analysis.View()
	// rightBox := m.help.View()

	switch m.state {
	// case idleState:
	// 	rightBox = m.help.View()
	case processingState:
		leftBox = m.analysis.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Render(leftBox,
		),
		m.statusbar.View(),
	)

	// return lipgloss.JoinVertical(lipgloss.Top,
	// 	lipgloss.NewStyle().Render(
	// 		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	// 	),
	// 	m.statusbar.View(),
	// )
}
