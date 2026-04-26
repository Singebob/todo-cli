package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"todo-cli/internal/domain"
)

type DetailModel struct {
	todo     domain.Todo
	viewport viewport.Model
	ready    bool
	width    int
	height   int
}

func NewDetailModel() DetailModel {
	return DetailModel{}
}

func (m DetailModel) SetTodo(todo domain.Todo, width, height int) DetailModel {
	m.todo = todo
	m.width = width
	m.height = height

	headerHeight := 6
	vpHeight := height - headerHeight
	if vpHeight < 1 {
		vpHeight = 1
	}

	m.viewport = viewport.New(width-4, vpHeight)
	m.viewport.SetContent(todo.Comment)
	m.ready = true
	return m
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	if !m.ready {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m DetailModel) View() string {
	var b strings.Builder

	status := "○ active"
	if m.todo.Done {
		status = "✓ completed"
		if m.todo.CompletedAt != nil {
			status += " " + m.todo.CompletedAt.Format("2006-01-02 15:04")
		}
	}

	b.WriteString(detailTitleStyle.Render(m.todo.Title))
	b.WriteString("\n")

	meta := fmt.Sprintf("Context: %s", m.todo.Context)
	if m.todo.Category != "" {
		meta += fmt.Sprintf("  |  Category: %s", m.todo.Category)
	}
	meta += fmt.Sprintf("  |  %s", status)
	meta += fmt.Sprintf("  |  Created: %s", m.todo.CreatedAt.Format("2006-01-02 15:04"))
	b.WriteString(detailMetaStyle.Render(meta))
	b.WriteString("\n\n")

	if m.todo.Comment != "" {
		b.WriteString(m.viewport.View())
	} else {
		b.WriteString(dimStyle.Render("(no comment)"))
	}

	return b.String()
}
