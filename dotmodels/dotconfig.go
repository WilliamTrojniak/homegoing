package dotmodels

import (
	"fmt"
	"slices"
	"strings"

	"github.com/WilliamTrojniak/homegoing/dotmanager"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DotConfigModel struct {
	Keys     keymap
	filepath string
	config   *dotmanager.DotConfig
	modules  []dotModuleModel
	tags     []dotModuleTagModel

	index    int
	tagIndex int
}

type keymap struct {
	Refresh key.Binding
	Link    key.Binding
	Unlink  key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
}

var keysDotConfig = keymap{
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh")),
	Link: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "link")),
	Unlink: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unlink")),
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "up")),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "down")),
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "right"),
	),
}

type getDotfilesConfigMsg struct {
	config *dotmanager.DotConfig
}

func NewDotConfigModel(filepath string) DotConfigModel {
	return DotConfigModel{filepath: filepath, Keys: keysDotConfig}
}

func (m *DotConfigModel) Load() tea.Cmd {
	return func() tea.Msg {
		config, err := dotmanager.LoadConfig(m.filepath)
		if err != nil {
			return err
		}
		return getDotfilesConfigMsg{config: config}
	}
}

func (m DotConfigModel) Init() tea.Cmd {
	return m.Load()
}

func (m DotConfigModel) Update(msg tea.Msg) (DotConfigModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO Handle resizing
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Refresh):
			return m, m.Load()
		case key.Matches(msg, m.Keys.Left):
			m.tagIndex = max(0, m.tagIndex-1)
			return m, m.tags[m.tagIndex].Init()
		case key.Matches(msg, m.Keys.Right):
			m.tagIndex = min(len(m.tags)-1, m.tagIndex+1)
			return m, m.tags[m.tagIndex].Init()
		}
	case getDotfilesConfigMsg:
		m, cmd := m.initConfig(msg)
		return m, cmd
	}
	var cmd tea.Cmd
	if len(m.tags) > 0 {
		m.tags[m.tagIndex], cmd = m.tags[m.tagIndex].Update(msg)
	}
	m, modCmd := m.updateModuleModels(msg)
	return m, tea.Batch(cmd, modCmd)
}

func (m DotConfigModel) initConfig(msg getDotfilesConfigMsg) (DotConfigModel, tea.Cmd) {
	m.config = msg.config
	m.tags = make([]dotModuleTagModel, 0)
	m.modules = make([]dotModuleModel, m.config.GetNumModules())
	tags := make(map[string][]*dotModuleModel)
	tags["All"] = []*dotModuleModel{}

	cmds := make([]tea.Cmd, m.config.GetNumModules())
	for i, moduleData := range m.config.GetModules() {
		m.modules[i] = NewDotModule(moduleData)
		module := &m.modules[i]
		cmds[i] = module.Init()
		tags["All"] = append(tags["All"], module)

		for _, tag := range moduleData.GetTags() {

			if modules, ok := tags[tag]; ok {
				tags[tag] = append(modules, module)
			} else {
				tags[tag] = []*dotModuleModel{module}
			}
		}

		// TODO: Init tags
	}

	for tag, modules := range tags {
		m.tags = append(m.tags, NewDotModuleTag(tag, modules, m.Keys))
	}

	slices.SortFunc(m.tags, func(t1, t2 dotModuleTagModel) int {
		if t1.tag < t2.tag {
			return -1
		}
		return 1
	})

	return m, tea.Batch(cmds...)

}

func (m DotConfigModel) updateModuleModels(msg tea.Msg) (DotConfigModel, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.modules))
	for i, moduleModel := range m.modules {
		m.modules[i], cmds[i] = moduleModel.Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m DotConfigModel) View() string {
	var b strings.Builder

	selectedTagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("4"))

	tagStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("8")).
		BorderRight(true)

	b.WriteString(" ")
	for i, tag := range m.tags {
		s := fmt.Sprintf(" %s (%v) ", tag.tag, len(tag.moduleModels))
		if i == m.tagIndex {
			b.WriteString(selectedTagStyle.Render(s))
		} else {
			b.WriteString(tagStyle.Render(s))
		}
		b.WriteString(" ")
	}
	b.WriteString("\n\n")

	if m.tagIndex < len(m.tags) {
		b.WriteString(m.tags[m.tagIndex].View())
	}

	return b.String()

}
