package models

import (
	"gohome/dotmanager"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

var lastId int;
var idMtx sync.Mutex;

func nextId() int {
  idMtx.Lock();
  defer idMtx.Unlock();
  lastId++;
  return lastId;
}

type Model struct {
  id int;
  *dotmanager.DotModule
  dotmanager.LinkStatus
}

type getStatusMsg struct {
  id int
  dotmanager.LinkStatus
}

type linkModuleMsg struct {
  id int
}

func NewModule(mod *dotmanager.DotModule) *Model { 
  return &Model{
    id: nextId(),
    DotModule: mod,
    LinkStatus: dotmanager.LINK_STATUS_UNKNOWN,
  }
}

func (m *Model) GetStatus() tea.Cmd {
  return func() tea.Msg {
    status, _ := m.DotModule.GetLinkStatus();
    return getStatusMsg{id: m.id, LinkStatus: status};
  }
}

func (m *Model) LinkModule(force bool) tea.Cmd {
  return func() tea.Msg {
    err := m.DotModule.LinkModule(force);
    if err != nil {
      return err;
    }
    return linkModuleMsg{id: m.id};
  }
}

func (m *Model) UnlinkModule() tea.Cmd {
  return func() tea.Msg {
    err := m.DotModule.UnlinkModule();
    if err != nil {
      return err;
    }
    return linkModuleMsg{id: m.id};
  }
}

func (m Model) Init() tea.Cmd {
  return m.GetStatus();
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case getStatusMsg:
    if msg.id != m.id {
      break;
    }
    m.LinkStatus = msg.LinkStatus;
    return m, nil;
  case linkModuleMsg:
    return m, m.GetStatus();
  }
  return m, nil;
}

func (m Model) View() string {
  return m.LinkStatus.String() + ": " + m.GetName() + "\n";
}


