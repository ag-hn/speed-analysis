package analysis

import (
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/ag-hn/speed-analysis/keys"
)

type processingState int

const (
	IdleState processingState = iota
	ProcessingState
)

type ProcessedItem struct {
	Name        string
	Addr     string
	Speed        string
	Lat   string
	Lng   string
}

type Model struct {
	processed             []ProcessedItem
	sub chan ProcessedItem
	count int

	Cursor                int
	min                   int
	max                   int
	height                int
	width                 int
	Disabled              bool
	keyMap                keys.KeyMap
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	selectedItemColor     lipgloss.AdaptiveColor
	unselectedItemColor   lipgloss.AdaptiveColor
	inactiveItemColor     lipgloss.AdaptiveColor
	err                   error
	State                 processingState

	enableLogging bool
}

func New() Model {
	return Model{
		Cursor: 0,
		min: 0,
		max: 0,
		sub: make(chan ProcessedItem),
		processed: []ProcessedItem{

		},
		Disabled:              false,
		keyMap:                keys.DefaultKeyMap(),
		StatusMessage:         "",
		StatusMessageLifetime: time.Second,
		selectedItemColor:     lipgloss.AdaptiveColor{Light: "212", Dark: "212"},
		unselectedItemColor:   lipgloss.AdaptiveColor{Light: "ffffff", Dark: "#000000"},
		inactiveItemColor:     lipgloss.AdaptiveColor{Light: "243", Dark: "243"},
	}
}
