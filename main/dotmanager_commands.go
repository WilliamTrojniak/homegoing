package main

import (
	"gohome/dotmanager"

	tea "github.com/charmbracelet/bubbletea"
)

type GetActiveSymLinksMsg struct { 
  Links []dotmanager.SymLink
}


func getActiveSymLinks(absLinkDestDirPath string, absLinkSrcDirPath string) tea.Cmd {
  return func() tea.Msg {
    
    symlinks, err := dotmanager.GetSymLinksInDir(absLinkDestDirPath, absLinkSrcDirPath);

    if err != nil {
      return ErrMsg{isFatal: false, err: err}
    }

    return GetActiveSymLinksMsg{Links: symlinks};
  }
}

