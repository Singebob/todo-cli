package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Enter      key.Binding
	Space      key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
	New        key.Binding
	Edit       key.Binding
	Delete     key.Binding
	Completed  key.Binding
	Filter     key.Binding
	Save       key.Binding
	Back       key.Binding
	Quit       key.Binding
	Yes        key.Binding
	No         key.Binding
	NewContext key.Binding
	DelContext key.Binding
}

var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k/↑", "up")),
	Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j/↓", "down")),
	Left:       key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("h/←", "left")),
	Right:      key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("l/→", "right")),
	Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
	Space:      key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle done")),
	Tab:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	ShiftTab:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
	New:        key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Completed:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "completed")),
	Filter:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	Save:       key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
	Back:       key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Yes:        key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yes")),
	No:         key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "no")),
	NewContext: key.NewBinding(key.WithKeys("C"), key.WithHelp("C", "new context")),
	DelContext: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete context")),
}
