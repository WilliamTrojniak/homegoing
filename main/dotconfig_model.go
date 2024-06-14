package main

import (
	"gohome/dotmanager"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type DotConfigModel struct {
  Keys keyMapDotConfig
  filepath string
  *dotmanager.DotConfig
  modules []DotModuleModel
  index int
}

type keyMapDotConfig struct {
  Refresh key.Binding
  Link key.Binding
  Unlink key.Binding
  Up key.Binding
  Down key.Binding
}

var keysDotConfig = keyMapDotConfig{
  Refresh: key.NewBinding(
    key.WithKeys("r"),
    key.WithHelp("r", "refresh")),
  Link: key.NewBinding(
    key.WithKeys("i",),
    key.WithHelp("i", "link")),
  Unlink: key.NewBinding(
    key.WithKeys("u"),
    key.WithHelp("u", "unlink")),
  Up: key.NewBinding(
    key.WithKeys("k",),
    key.WithHelp("k", "up")),
  Down: key.NewBinding(
    key.WithKeys("j",),
    key.WithHelp("j", "down")),
}

type GetDotfilesConfigMsg struct {
  config *dotmanager.DotConfig;
}

func NewDotConfigModel(filepath string) DotConfigModel {
  return DotConfigModel{filepath: filepath, Keys: keysDotConfig};
}

func (m *DotConfigModel) Load() tea.Cmd {
  return func() tea.Msg {
    config, err := dotmanager.LoadConfig(m.filepath);
    if err != nil {
      return err;
    }
    return GetDotfilesConfigMsg{config: config};
  }
}

func (m DotConfigModel) Init() tea.Cmd {
  return m.Load();
}

func (m DotConfigModel) Update(msg tea.Msg) (DotConfigModel, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    // TODO Handle resizing
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.Keys.Refresh):
      return m, m.Load();
    case key.Matches(msg, m.Keys.Up):
      m.index = max(m.index - 1, 0);
      return m, nil;
    case key.Matches(msg, m.Keys.Down):
      m.index = min(m.index + 1, len(m.modules) - 1);
      return m, nil;
    case key.Matches(msg, m.Keys.Link) && len(m.modules) > 0:
      return m, m.modules[m.index].LinkModule(false);
    case key.Matches(msg, m.Keys.Unlink) && len(m.modules) > 0:
      return m, m.modules[m.index].UnlinkModule();
    }
  case GetDotfilesConfigMsg:
    m, cmd := m.initConfig(msg);
    return m, cmd;
  }
  return m.updateModuleModels(msg);
}

func (m DotConfigModel) initConfig(msg GetDotfilesConfigMsg) (DotConfigModel, tea.Cmd) {
  m.DotConfig = msg.config;
  m.modules = make([]DotModuleModel, m.DotConfig.GetNumModules());
  cmds := make([]tea.Cmd, m.DotConfig.GetNumModules());
  for i, module := range m.DotConfig.GetModules() {
    m.modules[i] = NewDotModule(module);
    cmds[i] = m.modules[i].Init();
  }
  return m, tea.Batch(cmds...);

}

func (m DotConfigModel) updateModuleModels(msg tea.Msg) (DotConfigModel, tea.Cmd) {
  cmds := make([]tea.Cmd, len(m.modules));
  for i, moduleModel := range m.modules {
    m.modules[i], cmds[i] = moduleModel.Update(msg);
  }
  return m, tea.Batch(cmds...);

}

func (m DotConfigModel) View() string {
  var b strings.Builder;
  for _, module := range m.modules {
    b.WriteString(module.View());
  }
  return b.String();
}
