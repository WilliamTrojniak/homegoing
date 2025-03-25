package dotmodels

import (
	"strings"

	"github.com/WilliamTrojniak/homegoing/dotmanager"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
)

type dotModuleTagModel struct {
	id           int
	tag          string
	moduleModels []*dotModuleModel
	dotmanager.LinkStatus
	modelIndex int
	Keys       keymap
}

type initTagModelMsg struct {
	id int
}

func NewDotModuleTag(tag string, mod []*dotModuleModel, keys keymap) dotModuleTagModel {
	return dotModuleTagModel{
		id:           nextId(),
		tag:          tag,
		moduleModels: mod,
		LinkStatus:   dotmanager.LINK_STATUS_UNKNOWN,
		modelIndex:   0,
		Keys:         keys,
	}
}

func (m dotModuleTagModel) Init() tea.Cmd {
	return func() tea.Msg {
		return initTagModelMsg{id: m.id}
	}
}

func (m dotModuleTagModel) Update(msg tea.Msg) (dotModuleTagModel, tea.Cmd) {
	switch msg := msg.(type) {
	case initTagModelMsg:
		if msg.id != m.id {
			break
		}
		m.modelIndex = 0
		cmds := make([]tea.Cmd, len(m.moduleModels))
		for i, module := range m.moduleModels {
			cmds[i] = module.Init()
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			m.modelIndex = max(0, m.modelIndex-1)
			return m, nil
		case key.Matches(msg, m.Keys.Down):
			m.modelIndex = min(len(m.moduleModels)-1, m.modelIndex+1)
			return m, nil
		case key.Matches(msg, m.Keys.Link) && len(m.moduleModels) > 0:
			return m, m.moduleModels[m.modelIndex].LinkModule(false)
		case key.Matches(msg, m.Keys.Unlink) && len(m.moduleModels) > 0:
			return m, m.moduleModels[m.modelIndex].UnlinkModule()
		}
	}
	return m, nil
}

func (m dotModuleTagModel) View() string {
	var b strings.Builder
	l := list.New()
	for _, module := range m.moduleModels {
		l.Item(module.View())
	}

	l.Enumerator(func(l list.Items, i int) string {
		// if i == m.index {
		if i == m.modelIndex {
			return ">"
		}
		return ""
	})

	b.WriteString(l.String())
	b.WriteString("\n\n")
	return b.String()
}
