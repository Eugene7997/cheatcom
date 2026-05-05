package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Eugene7997/cheatcom/internal/action"
	"github.com/Eugene7997/cheatcom/internal/store"
)

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))

type Result struct {
	Cheat     store.Cheat
	Action    action.Action
	Cancelled bool
}

type model struct {
	list          list.Model
	result        Result
	cancelled     bool
	chosen        bool
	defaultAction action.Action
}

func newModel(cheats []store.Cheat, defaultAction action.Action) model {
	items := make([]list.Item, len(cheats))
	for i, c := range cheats {
		items[i] = cheatItem{c: c}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Cheats"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetStatusBarItemName("cheat", "cheats")

	return model{list: l, defaultAction: defaultAction}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-2)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, keys.Quit):
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, keys.Enter):
			item, ok := m.list.SelectedItem().(cheatItem)
			if !ok {
				return m, nil
			}
			m.result = Result{Cheat: item.c, Action: m.defaultAction}
			m.chosen = true
			return m, tea.Quit

		case key.Matches(msg, keys.Run):
			item, ok := m.list.SelectedItem().(cheatItem)
			if !ok {
				return m, nil
			}
			m.result = Result{Cheat: item.c, Action: action.Run}
			m.chosen = true
			return m, tea.Quit

		case key.Matches(msg, keys.Print):
			item, ok := m.list.SelectedItem().(cheatItem)
			if !ok {
				return m, nil
			}
			m.result = Result{Cheat: item.c, Action: action.Print}
			m.chosen = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) defaultLabel() string {
	switch m.defaultAction {
	case action.Run:
		return "run"
	case action.Print:
		return "print"
	default:
		return "copy"
	}
}

func (m model) View() string {
	if len(m.list.Items()) == 0 {
		return lipgloss.NewStyle().Margin(2).Render(
			"No cheats yet.\n\nAdd one with:  chc add \"<command>\" -d \"<description>\"",
		)
	}
	return m.list.View() + "\n  enter:" + m.defaultLabel() + "  r:run  p:print  /:filter  q:quit"
}

func Run(cheats []store.Cheat, defaultAction action.Action) (Result, error) {
	m := newModel(cheats, defaultAction)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return Result{}, err
	}
	fm := final.(model)
	if fm.cancelled || !fm.chosen {
		return Result{Cancelled: true}, nil
	}
	return fm.result, nil
}
