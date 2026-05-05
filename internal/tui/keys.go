package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Enter  key.Binding
	Run    key.Binding
	Print  key.Binding
	Quit   key.Binding
	Filter key.Binding
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Run: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "run"),
	),
	Print: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "print"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
}
