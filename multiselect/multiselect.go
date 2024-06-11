package multiselect

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
  Up key.Binding
  Down key.Binding
  Toggle key.Binding
  Confirm key.Binding
  Cancel key.Binding
}

var keys = keyMap {
  Up: key.NewBinding(
    key.WithKeys("up", "k"),
    key.WithHelp("/k", "move up")),
  Down: key.NewBinding(
    key.WithKeys("down", "j"),
    key.WithHelp("/j", "move down")),
  Toggle: key.NewBinding(
    key.WithKeys("space"),
    key.WithHelp("󱁐", "toggle")),
  Confirm: key.NewBinding(
    key.WithKeys("y"),
    key.WithHelp("y", "confirm")),
  Cancel: key.NewBinding(
    key.WithKeys("x"),
    key.WithHelp("x", "cancel")),
}

type Option interface {
  GetValue() string
  GetId() string
}


type MultiSelectModel[T Option] struct {
  Keys keyMap
  Options []T
}

func NewMultiSelect[T Option](options []T) *MultiSelectModel[T] {
  return &MultiSelectModel[T]{Keys: keys, Options: options};
}


func (m MultiSelectModel[T]) Init() tea.Cmd {
  return nil;
}

func (m MultiSelectModel[T]) Update(tea.Msg) (tea.Model, tea.Cmd) {
  return m, nil;
}

func (m MultiSelectModel[T]) View() string {
  var builder strings.Builder;
  for _, option := range m.Options {
    builder.WriteString(fmt.Sprintf("[ ] %v\n", option.GetValue()));
  }

  return builder.String();
}
