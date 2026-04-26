package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 2)

	tabInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("238")).
			Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
			MarginBottom(1)

	tabAdd = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 1)
)

type TabBarModel struct {
	contexts []string
	cursor   int
	width    int
}

func NewTabBarModel() TabBarModel {
	return TabBarModel{}
}

func (m TabBarModel) SelectedContext() string {
	if len(m.contexts) == 0 {
		return ""
	}
	if m.cursor >= len(m.contexts) {
		return m.contexts[0]
	}
	return m.contexts[m.cursor]
}

func (m TabBarModel) Next() TabBarModel {
	if len(m.contexts) > 0 {
		m.cursor = (m.cursor + 1) % len(m.contexts)
	}
	return m
}

func (m TabBarModel) Prev() TabBarModel {
	if len(m.contexts) > 0 {
		m.cursor = (m.cursor - 1 + len(m.contexts)) % len(m.contexts)
	}
	return m
}

func (m TabBarModel) SelectByName(name string) TabBarModel {
	for i, c := range m.contexts {
		if c == name {
			m.cursor = i
			return m
		}
	}
	return m
}

func (m TabBarModel) View() string {
	if len(m.contexts) == 0 {
		return tabBarStyle.Render(dimStyle.Render("No contexts — press C to create one"))
	}

	var tabs []string
	for i, ctx := range m.contexts {
		if i == m.cursor {
			tabs = append(tabs, tabActive.Render(ctx))
		} else {
			tabs = append(tabs, tabInactive.Render(ctx))
		}
	}

	row := strings.Join(tabs, " ")
	row += " " + tabAdd.Render("[C:new ctx]")

	return tabBarStyle.Width(m.width).Render(row)
}
