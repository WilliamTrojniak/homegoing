package main

import (
	"gohome/dotmanager"

	tea "github.com/charmbracelet/bubbletea"
)

type GetDotfilesConfigMsg struct {
  config *dotmanager.DotConfig;
}

type LinkDotModuleMsg struct {
  module *dotmanager.DotModule
}

type GetActiveSymLinksMsg struct { 
  links []dotmanager.SymLink
}

func getDotfilesConfig(path string) tea.Cmd {
  return func() tea.Msg {
    config, err := dotmanager.LoadConfig(path);

    if err != nil {
      return ErrMsg{isFatal: false, err: err};
    }

    return GetDotfilesConfigMsg{config: config};
  }
}

func linkModule(mod *dotmanager.DotModule) tea.Cmd {
  return func() tea.Msg {
    if err := mod.LinkModule(true); err != nil {
      return ErrMsg{isFatal: false, err: err};
    }
    return LinkDotModuleMsg{module: mod};
  }
}

func unlinkModule(mod *dotmanager.DotModule) tea.Cmd {
  return func() tea.Msg {
    if err := mod.UnlinkModule(); err != nil {
      return ErrMsg{isFatal: false, err: err};
    }
    return LinkDotModuleMsg{module: mod};
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

