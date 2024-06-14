package main

import (
	"gohome/dotmanager"
	"strings"
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

type DotModuleModel struct {
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

func NewDotModule(mod *dotmanager.DotModule) DotModuleModel { 
  return DotModuleModel{
    id: nextId(),
    DotModule: mod,
    LinkStatus: dotmanager.LINK_STATUS_UNKNOWN,
  }
}

func (m *DotModuleModel) GetStatus() tea.Cmd {
  return func() tea.Msg {
    status, _ := m.DotModule.GetLinkStatus();
    return getStatusMsg{id: m.id, LinkStatus: status};
  }
}

func (m *DotModuleModel) LinkModule(force bool) tea.Cmd {
  return func() tea.Msg {
    err := m.DotModule.LinkModule(force);
    if err != nil {
      return err;
    }
    return linkModuleMsg{id: m.id};
  }
}

func (m *DotModuleModel) UnlinkModule() tea.Cmd {
  return func() tea.Msg {
    err := m.DotModule.UnlinkModule();
    if err != nil {
      return err;
    }
    return linkModuleMsg{id: m.id};
  }
}

func (m DotModuleModel) Init() tea.Cmd {
  return m.GetStatus();
}

func (m DotModuleModel) Update(msg tea.Msg) (DotModuleModel, tea.Cmd) {
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

func (m DotModuleModel) View() string {
  var s strings.Builder;
  switch m.LinkStatus {
  case dotmanager.LINK_STATUS_UNLINKED:
    s.WriteString("[ ] ");
  case dotmanager.LINK_STATUS_LINKED:
    s.WriteString("[󰄬] ");
  case dotmanager.LINK_STATUS_EXISTS_CONFLICT, dotmanager.LINK_STATUS_TARGET_CONFLICT:
    s.WriteString("[!] ");
  case dotmanager.LINK_STATUS_UNKNOWN:
    s.WriteString("[?] ");
  }

  s.WriteString(m.GetName());
  s.WriteString("\n");

  return s.String();
}


