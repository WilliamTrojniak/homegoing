package dotmodels

import (
	"fmt"
	"strings"
	"sync"

	"github.com/WilliamTrojniak/homegoing/dotmanager"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastId int
var idMtx sync.Mutex

func nextId() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastId++
	return lastId
}

type dotModuleModel struct {
	id int
	dotmanager.DotModule
	dotmanager.LinkStatus
}

type getStatusMsg struct {
	id int
	dotmanager.LinkStatus
}

type linkModuleMsg struct {
	id int
}

func NewDotModule(mod dotmanager.DotModule) dotModuleModel {
	return dotModuleModel{
		id:         nextId(),
		DotModule:  mod,
		LinkStatus: dotmanager.LINK_STATUS_UNKNOWN,
	}
}

func (m *dotModuleModel) GetStatus() tea.Cmd {
	return func() tea.Msg {
		status, _ := m.DotModule.GetLinkStatus()
		return getStatusMsg{id: m.id, LinkStatus: status}
	}
}

func (m *dotModuleModel) LinkModule(force bool) tea.Cmd {
	return func() tea.Msg {
		err := m.DotModule.LinkModule(force)
		if err != nil {
			return err
		}
		return linkModuleMsg{id: m.id}
	}
}

func (m *dotModuleModel) UnlinkModule() tea.Cmd {
	return func() tea.Msg {
		err := m.DotModule.UnlinkModule()
		if err != nil {
			return err
		}
		return linkModuleMsg{id: m.id}
	}
}

func (m dotModuleModel) Init() tea.Cmd {
	return m.GetStatus()
}

func (m dotModuleModel) Update(msg tea.Msg) (dotModuleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case getStatusMsg:
		if msg.id != m.id {
			break
		}
		m.LinkStatus = msg.LinkStatus
		return m, nil
	case linkModuleMsg:
		return m, m.GetStatus()
	}
	return m, nil
}

func (m dotModuleModel) View() string {
	var s strings.Builder
	switch m.LinkStatus {
	case dotmanager.LINK_STATUS_UNLINKED:
		s.WriteString("[ ] ")
	case dotmanager.LINK_STATUS_LINKED:
		s.WriteString("[󰄬] ")
	case dotmanager.LINK_STATUS_EXISTS_CONFLICT, dotmanager.LINK_STATUS_TARGET_CONFLICT:
		s.WriteString("[!] ")
	case dotmanager.LINK_STATUS_UNKNOWN:
		s.WriteString("[?] ")
	}

	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).Faint(true)

	s.WriteString(m.GetName())
	s.WriteString(tagStyle.Render(fmt.Sprintf(" [%s]", strings.Join(m.GetTags(), ", "))))

	return s.String()
}
