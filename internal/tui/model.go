package tui

import (
	"github.com/ag-hn/speed-analysis/analysis"
	"github.com/ag-hn/speed-analysis/help"
	"github.com/ag-hn/speed-analysis/keys"
)

type sessionState int

const (
	idleState sessionState = iota
	showHelpState
	processingState
)

type Config struct {
	EnableLogging  bool
}

type model struct {
	help                  help.Model
	analysis  			analysis.Model
	state                 sessionState
	keyMap                keys.KeyMap
}

// New creates a new instance of the UI.
func New(cfg Config) *model {
	analysis := analysis.New()

	defaultKeyMap := keys.DefaultKeyMap()

	helpModel := help.New(
		"Help",
		help.TitleColor{
		},
		[]help.Entry{
			{Key: defaultKeyMap.ForceQuit.Help().Key, Description: defaultKeyMap.ForceQuit.Help().Desc},
			{Key: defaultKeyMap.Quit.Help().Key, Description: defaultKeyMap.Quit.Help().Desc},
			// {Key: defaultKeyMap.Process.Help().Key, Description: defaultKeyMap.Process.Help().Desc},
			{Key: defaultKeyMap.CopyPathToClipboard.Help().Key, Description: defaultKeyMap.CopyPathToClipboard.Help().Desc},	
			{Key: defaultKeyMap.Down.Help().Key, Description: defaultKeyMap.Down.Help().Desc},
			{Key: defaultKeyMap.Up.Help().Key, Description: defaultKeyMap.Up.Help().Desc},	
		},
	)
	helpModel.SetViewportDisabled(true)

	return &model{
		help:                  helpModel,
		analysis: analysis,
		state: idleState,
		keyMap: defaultKeyMap,
	}
}
