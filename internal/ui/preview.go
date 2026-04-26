package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"todo-cli/internal/domain"
)

var (
	previewBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			BorderLeft(true).
			Padding(0, 1)
)

type PreviewModel struct {
	todo     *domain.Todo
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

func NewPreviewModel() PreviewModel {
	return PreviewModel{}
}

func (m PreviewModel) SetTodo(todo *domain.Todo, width, height int) PreviewModel {
	m.todo = todo
	m.width = width
	m.height = height

	innerW := width - 4
	if innerW < 10 {
		innerW = 10
	}

	headerLines := 5
	vpH := height - headerLines
	if vpH < 1 {
		vpH = 1
	}

	m.viewport = viewport.New(innerW, vpH)
	if todo != nil && todo.Comment != "" {
		m.viewport.SetContent(todo.Comment)
	} else {
		m.viewport.SetContent("")
	}
	m.ready = true
	return m
}

func (m PreviewModel) Update(msg tea.Msg) (PreviewModel, tea.Cmd) {
	if !m.ready || m.todo == nil {
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m PreviewModel) View() string {
	if m.todo == nil {
		content := dimStyle.Render("No todo selected")
		return previewBorder.Width(m.width).Height(m.height).Render(content)
	}

	var b strings.Builder
	b.WriteString(detailTitleStyle.Render(m.todo.Title))
	b.WriteString("\n")

	meta := fmt.Sprintf("%s", m.todo.Context)
	if m.todo.Category != "" {
		meta += " | " + m.todo.Category
	}
	meta += " | " + formatAge(m.todo.CreatedAt)
	if m.todo.Done && m.todo.CompletedAt != nil {
		meta += " | done " + formatAge(*m.todo.CompletedAt)
	}
	b.WriteString(detailMetaStyle.Render(meta))
	b.WriteString("\n\n")

	if m.todo.Comment != "" {
		b.WriteString(m.viewport.View())
	} else {
		b.WriteString(dimStyle.Render("(no comment)"))
	}

	return previewBorder.Width(m.width).Height(m.height).Render(b.String())
}
