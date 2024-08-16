package analysis

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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

func (m Model) waitForProcessedItem() tea.Cmd {
	return func() tea.Msg {
		return waitForProcessedItemMsg{
			processed: <-m.sub,
		}
	}
}

func (m Model) ListenForProcessedItem() tea.Cmd {
	return func() tea.Msg {
		paths, err := ListProcessFilePaths()
		if err != nil {
			panic(err)
		}

		for _, path := range paths {
			// time.Sleep(time.Millisecond * 100)
			item, err := ProcessFilePath(path)
			if err != nil {
				m.err = err
				return errorMsg(err.Error())
			}

			for _, i := range item {
				m.sub <- i
			}
		}

		return listenForProcessedItemMsg("Done!")
	}
}

func ProcessedItemToString(p ProcessedItem) string {
	return fmt.Sprintf("Name: %s | Addr: %s | (Lat,Lng): (%s,%s) | Speed: %s mph", p.Name, p.Addr, p.Lat, p.Lng, p.Speed)
}

// copyToClipboardCmd copies the provided string to the clipboard.
func copyToClipboardCmd(p ProcessedItem) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.WriteAll(ProcessedItemToString(p))

		if err != nil {
			return errorMsg(err.Error())
		}

		return copyToClipboardMsg(
			fmt.Sprintf("%s %s %s", "Successfully copied", p.Name, "to clipboard"),
		)
	}
}

func openEditorCmd(p ProcessedItem) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(
			"openEditorCmd+Not implemented",
		)
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "code"
		}

		file := path.Join(OUTPUT_DATA_FILE, p.Addr)
		c := exec.Command(editor, file)

		return tea.ExecProcess(c, func(err error) tea.Msg {
			return editorFinishedMsg{err}
		})
	}
}
