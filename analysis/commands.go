package analysis

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type errorMsg string
type statusMessageTimeoutMsg struct{}
type editorFinishedMsg struct{ err error }
type getProcessedFilesMsg struct {
	processed []ProcessedItem
}
type waitForProcessedItemMsg struct {
	processed ProcessedItem
}
type listenForProcessedItemMsg string
type copyToClipboardMsg string

// NewStatusMessageCmd sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *Model) NewStatusMessageCmd(s string) tea.Cmd {
	m.StatusMessage = s

	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

// GetDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) GetDirectoryListingCmd() tea.Cmd {
	return func() tea.Msg {

		return getProcessedFilesMsg{
			processed: []ProcessedItem{
				{
					Name:  "Fake",
					Addr:  "asldfjeiojwaklcdjkljersa",
					Speed: "100",
					Lat:   "0",
					Lng:   "0",
				},
				{
					Name:  "Fake 2",
					Addr:  "asldfjeiojwaklcdjkljersa",
					Speed: "100",
					Lat:   "0",
					Lng:   "0",
				},
			},
		}
	}
}

func (m Model) waitForProcessedItem() tea.Cmd {
	return func() tea.Msg {
		return waitForProcessedItemMsg{
			processed: <-m.sub,
		}
	}
}

func (m Model) ListenForProcessedItem() tea.Cmd {
	return func() tea.Msg {
		for i := range 10 {
			time.Sleep(time.Millisecond * 100)
			m.sub <- ProcessedItem{
				Name:  fmt.Sprintf("%d", i),
				Addr:  "asodijcpsklJopasdfjelsadfc",
				Speed: "0",
				Lat:   "0",
				Lng:   "0",
			}
		}

		return listenForProcessedItemMsg("Done!")
	}
}

func processedItemToString(p ProcessedItem) string {
	return fmt.Sprintf("Name: %s | Addr: %s | (Lat,Lng): (%s,%s) | Speed: %s mph", p.Name, p.Addr, p.Lat, p.Lng, p.Speed)
}

// copyToClipboardCmd copies the provided string to the clipboard.
func copyToClipboardCmd(p ProcessedItem) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(processedItemToString(p))
		fmt.Println("printing", err)
		if err != nil {
			return errorMsg(err.Error())
		}

		return copyToClipboardMsg(
			fmt.Sprintf("%s %s %s", "Successfully copied", p.Name, "to clipboard"),
		)
	}
}
