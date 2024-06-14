package main

import (
	"gohome/dotmanager"

	tea "github.com/charmbracelet/bubbletea"
)

type GetDotfilesConfigMsg struct {
  config *dotmanager.DotConfig;
}


func getDotfilesConfig(path string) tea.Cmd {
  return func() tea.Msg {
    config, err := dotmanager.LoadConfig(path);

    if err != nil {
      return ErrMsg{isFatal: false, error: err};
    }

    return GetDotfilesConfigMsg{config: config};
  }
}
