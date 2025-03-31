package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Copy     key.Binding
	Back     key.Binding
	Quit     key.Binding
	Help     key.Binding
	Create   key.Binding
	Confirm  key.Binding
	Cancel   key.Binding
	Search   key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
}

func New() KeyMap {
	return KeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Select:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/confirm")),
		Copy:     key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy XML")),
		Back:     key.NewBinding(key.WithKeys("left", "esc"), key.WithHelp("←/esc", "back/cancel")),
		Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Create:   key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new prompt")),
		Confirm:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		Cancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "previous field")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Help, k.Create, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Search, k.Select},
		{k.Copy, k.Back},
		{k.Create},
		{k.Help, k.Quit},
	}
}

func Matches(msg tea.KeyMsg, binding key.Binding) bool {
	return key.Matches(msg, binding)
}

func MatchesKeyType(msg tea.KeyMsg, keyType tea.KeyType) bool {
	return msg.Type == keyType
}
