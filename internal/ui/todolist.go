package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"todo-cli/internal/domain"
)

type TodoListModel struct {
	todos        []domain.Todo
	cursor       int
	scrollOffset int
	focused      bool
	width        int
	height       int

	showCompleted  bool
	categoryFilter string
}

func NewTodoListModel() TodoListModel {
	return TodoListModel{}
}

func (m TodoListModel) Update(msg tea.Msg) (TodoListModel, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scrollOffset {
					m.scrollOffset = m.cursor
				}
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.todos)-1 {
				m.cursor++
				visible := m.visibleLines()
				if m.cursor >= m.scrollOffset+visible {
					m.scrollOffset = m.cursor - visible + 1
				}
			}
		}
	}
	return m, nil
}

func (m TodoListModel) visibleLines() int {
	v := m.height - 4
	if v < 1 {
		return 1
	}
	return v
}

func (m TodoListModel) SelectedTodo() (domain.Todo, bool) {
	if len(m.todos) == 0 {
		return domain.Todo{}, false
	}
	if m.cursor >= len(m.todos) {
		return m.todos[0], true
	}
	return m.todos[m.cursor], true
}

func (m TodoListModel) View() string {
	var b strings.Builder

	title := "Todos"
	if m.showCompleted {
		title += " (completed)"
	}
	if m.categoryFilter != "" {
		title += fmt.Sprintf(" [%s]", m.categoryFilter)
	}
	b.WriteString(detailTitleStyle.Render(title))
	b.WriteString("\n\n")

	if len(m.todos) == 0 {
		if m.showCompleted {
			b.WriteString(dimStyle.Render("No completed todos"))
		} else {
			b.WriteString(dimStyle.Render("No todos yet — press n to create one"))
		}
		return listPaneStyle.Width(m.width).Height(m.height).Render(b.String())
	}

	visible := m.visibleLines()
	end := m.scrollOffset + visible
	if end > len(m.todos) {
		end = len(m.todos)
	}

	for i := m.scrollOffset; i < end; i++ {
		t := m.todos[i]
		check := "[ ]"
		if t.Done {
			check = "[✓]"
		}

		age := dimStyle.Render(formatAge(t.CreatedAt))
		if t.Done && t.CompletedAt != nil {
			age = dimStyle.Render("✓ " + formatAge(*t.CompletedAt))
		}

		line := fmt.Sprintf("%s %s", check, t.Title)
		if t.Category != "" {
			line += " " + categoryTag.Render("#"+t.Category)
		}
		line += "  " + age

		if i == m.cursor {
			if t.Done {
				b.WriteString(todoSelected.Render(todoDone.Render(line)))
			} else {
				b.WriteString(todoSelected.Render(line))
			}
		} else {
			if t.Done {
				b.WriteString(todoDone.Render(line))
			} else {
				b.WriteString(todoNormal.Render(line))
			}
		}
		b.WriteString("\n")
	}

	if len(m.todos) > visible {
		b.WriteString(dimStyle.Render(fmt.Sprintf("\n %d/%d", m.cursor+1, len(m.todos))))
	}

	return listPaneStyle.Width(m.width).Height(m.height).Render(b.String())
}
