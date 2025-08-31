package tui

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	quit key.Binding
	hide key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("ctrl+c/q/esc", "quit"),
		),
		hide: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "hide"),
		),
	}
}
