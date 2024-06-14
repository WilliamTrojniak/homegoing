package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
  Help key.Binding
  Quit key.Binding
}

var keys = keyMap{
  Quit: key.NewBinding(
    key.WithKeys("q", "ctrl+c"),
    key.WithHelp("q", "quit")),
  Help: key.NewBinding(
    key.WithKeys("?"),
    key.WithHelp("?", "show help")),
}

func (k keyMap) ShortHelp() []key.Binding {
  return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
  return [][]key.Binding{{k.Help},{k.Quit}};
}


type app struct {
  help help.Model
  keys keyMap
  config DotConfigModel

  error
  isQuitting bool
}

func newApp(configFilePath string) *app {
  help := help.New();
  help.ShowAll = true;

  return &app{
    help: help, 
    keys: keys, 
    config: NewDotConfigModel(configFilePath),
  };
}

func (m app) Init() tea.Cmd {
  return m.config.Init();
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

  switch msg := msg.(type) {
  case error:
   // TODO: Log other errors
    m.error = msg;
    return m, nil;
  case tea.WindowSizeMsg:
    // TODO: Handle resizing
  case tea.KeyMsg:
    if m.error != nil {
      m.error = nil;
      return m, nil;
    }
    switch {
    case key.Matches(msg, m.keys.Quit):
      m.isQuitting = true;
      return m, tea.Quit
    case key.Matches(msg, m.keys.Help):
    // TODO Implement
    }
  }
  var cmd tea.Cmd;
  m.config, cmd = m.config.Update(msg);

  return m, cmd;
}

func (m app) View() string {
  if m.isQuitting {
    if m.error != nil {
      return fmt.Sprintf("A fatal error occurred: %v", m.error);
    }
    return "";
  }
  if m.error != nil {
    return fmt.Sprintf("%v", m.error);
  }
  var s string;
  s = m.config.View();
  return s;
}
