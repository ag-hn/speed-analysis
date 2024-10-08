package analysis

// SetDisabled sets if the bubble is currently active.
func (m *Model) SetDisabled(disabled bool) {
	m.Disabled = disabled
}

func (m *Model) SetDebugging(debugging bool) {
	m.enableLogging = debugging
}

// GetSelectedItem returns the currently selected file/dir.
func (m Model) GetSelectedItem() ProcessedItem {
	if len(m.processed) > 0 {
		return m.processed[m.Cursor]
	}

	return ProcessedItem{}
}

// SetSize Sets the size of the filetree.
func (m *Model) SetSize(width, height int) {
	m.height = height
	m.width = width
	m.max = m.height - 1
}

// GetTotalItems returns total number of tree items.
func (m Model) GetTotalItems() int {
	return len(m.processed)
}
