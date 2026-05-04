package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Eugene7997/cheatcom/internal/store"
)

var labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("241"))
var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

type editModel struct {
	inputs    []textinput.Model
	focused   int
	submitted bool
	cancelled bool
}

const (
	fieldCommand = iota
	fieldDescription
	fieldTags
)

func newEditModel(c store.Cheat) editModel {
	inputs := make([]textinput.Model, 3)
	for i := range inputs {
		t := textinput.New()
		t.CharLimit = 512
		inputs[i] = t
	}

	inputs[fieldCommand].Placeholder = "command"
	inputs[fieldCommand].SetValue(c.Command)
	inputs[fieldCommand].Focus()

	inputs[fieldDescription].Placeholder = "description"
	inputs[fieldDescription].SetValue(c.Description)

	inputs[fieldTags].Placeholder = "tag1,tag2"
	inputs[fieldTags].SetValue(strings.Join(c.Tags, ","))

	return editModel{inputs: inputs, focused: 0}
}

func (m editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		case "esc":
			m.cancelled = true
			return m, tea.Quit
		case "tab", "down":
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % len(m.inputs)
			m.inputs[m.focused].Focus()
		case "shift+tab", "up":
			m.inputs[m.focused].Blur()
			m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.focused].Focus()
		case "enter":
			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				m.inputs[m.focused].Focus()
			} else {
				m.submitted = true
				return m, tea.Quit
			}
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m editModel) View() string {
	labels := []string{"Command:     ", "Description: ", "Tags:        "}
	var b strings.Builder
	b.WriteString("\n  Edit cheat\n\n")
	for i, inp := range m.inputs {
		label := labelStyle.Render(labels[i])
		if i == m.focused {
			label = focusedStyle.Render(labels[i])
		}
		b.WriteString("  " + label + inp.View() + "\n\n")
	}
	b.WriteString("  tab/↑↓ navigate  enter next/submit  esc cancel\n")
	return b.String()
}

type EditResult struct {
	Command     string
	Description string
	Tags        []string
	Submitted   bool
}

func RunEdit(c store.Cheat) (EditResult, error) {
	m := newEditModel(c)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return EditResult{}, err
	}
	fm := final.(editModel)
	if fm.cancelled || !fm.submitted {
		return EditResult{}, nil
	}

	rawTags := strings.Split(fm.inputs[fieldTags].Value(), ",")
	tags := normalizeTags(rawTags)

	return EditResult{
		Command:     fm.inputs[fieldCommand].Value(),
		Description: fm.inputs[fieldDescription].Value(),
		Tags:        tags,
		Submitted:   true,
	}, nil
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(strings.ToLower(t))
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
