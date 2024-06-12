package main

import (
	"gohome/dotmanager"

	tea "github.com/charmbracelet/bubbletea"
)

type GetDotfilesConfigMsg struct {
  config dotmanager.DotConfig;
}

type GetActiveSymLinksMsg struct { 
  links []dotmanager.SymLink
}

func getDotfilesConfig(path string) tea.Cmd {
  return func() tea.Msg {
    config, err := dotmanager.ReadConfig(path);

    if err != nil {
      return ErrMsg{isFatal: false, err: err};
    }

    return GetDotfilesConfigMsg{config: config};
  }
}

func getActiveSymLinks(absLinkDestDirPath string, absLinkSrcDirPath string) tea.Cmd {
  return func() tea.Msg {
    
    symlinks, err := dotmanager.GetSymLinksInDir(absLinkDestDirPath, absLinkSrcDirPath);

    if err != nil {
      return ErrMsg{isFatal: false, err: err}
    }

    return GetActiveSymLinksMsg{links: symlinks};
  }
}

