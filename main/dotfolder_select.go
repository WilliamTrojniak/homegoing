package main

import (
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/filepicker"
)

type DotfolderSelectModel struct {
  fp filepicker.Model
}

func (m DotfolderSelectModel) Init() tea.Cmd {
  return nil;
}

func (m DotfolderSelectModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
  return m, nil;
}

func (m DotfolderSelectModel) View() string {
  return "Filepicker";
}
