package tui

import (
	"fmt"

	"github.com/ag-hn/speed-analysis/analysis"
	"github.com/ag-hn/speed-analysis/polish"
	"github.com/charmbracelet/lipgloss"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (m *model) disableAllViewports() {
	m.help.SetViewportDisabled(true)
}

func (m *model) resetViewports() {
	m.help.GotoTop()
}

func (m *model) updateStatusBar() {
	if m.analysis.GetSelectedItem().Name != "" {
		statusMessage :=
			lipgloss.NewStyle().
					Padding(0, 1).
					Foreground(polish.Colors.Yellow500).
					Render(analysis.ProcessedItemToString(m.analysis.GetSelectedItem()))

		if m.analysis.StatusMessage != "" {
			statusMessage = m.analysis.StatusMessage
		}

		if m.statusMessage != "" {
			statusMessage = m.statusMessage
		}

		m.statusbar.SetContent(
			m.analysis.GetSelectedItem().Name,
			statusMessage,
			fmt.Sprintf("%d/%d", m.analysis.Cursor+1, m.analysis.GetTotalItems()),
			fmt.Sprintf(m.analysis.GetSelectedItem().Speed),
		)
	} else {
		statusMessage := "No items processed"

		m.statusbar.SetContent(
			"N/A",
			statusMessage,
			fmt.Sprintf("%d/%d", 0, 0),
			"FM",
		)
	}
}
