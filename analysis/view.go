package analysis

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var fileList strings.Builder

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	if m.State == ProcessingState {
		return lipgloss.NewStyle().
			Width(m.width).
			Render("Processing...", fmt.Sprintf(" (%d/???)", m.count)) + "\n"
	}

	for i, file := range m.processed {
		if i < m.min || i > m.max {
			continue
		}

		switch {
		case m.Disabled:
			fallthrough

		case i == m.Cursor && !m.Disabled:
			textColor := m.inactiveItemColor

			if i == m.Cursor && !m.Disabled {
				textColor = m.selectedItemColor
			}

			fileList.WriteString(
				renderProcessedItem(file, textColor),
			)

		case i != m.Cursor && !m.Disabled:
			textColor := m.unselectedItemColor

			fileList.WriteString(
				renderProcessedItem(file, textColor),
			)
		}
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Render(fileList.String())
}

func renderProcessedItem(p ProcessedItem, textColor lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		Render(processedItemToString(p)) + "\n"
}
