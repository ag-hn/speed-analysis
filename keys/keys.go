package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down                key.Binding
	Up                  key.Binding
	ForceQuit           key.Binding
	Quit                key.Binding
	Process key.Binding
	CopyPathToClipboard key.Binding	
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ForceQuit:           key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Force Quit")),
		Quit:                key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit when not performing action")),
		Process:             key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "Begin processing average speeds")),
		Down:                key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "Go down")),
		Up:                  key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "Go up")),	
		CopyPathToClipboard: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "Copy path to clipboard")),	
	}
}
