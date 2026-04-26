package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"todo-cli/internal/domain"
)

const (
	fieldTitle = iota
	fieldContext
	fieldCategory
	fieldComment
	fieldCount
)

type FormModel struct {
	title    textinput.Model
	category textinput.Model
	comment  textarea.Model

	contexts     []string
	contextIndex int

	focusIndex int
	editing    bool
	editID     uuid.UUID

	width  int
	height int
	err    string
}

func NewFormModel() FormModel {
	ti := textinput.New()
	ti.Placeholder = "Todo title"
	ti.CharLimit = 120

	cati := textinput.New()
	cati.Placeholder = "Category (optional)"
	cati.CharLimit = 40

	cmt := textarea.New()
	cmt.Placeholder = "Comment — links, notes, context…"
	cmt.CharLimit = 2000
	cmt.SetHeight(6)

	return FormModel{
		title:    ti,
		category: cati,
		comment:  cmt,
	}
}

func (m FormModel) SetForCreate(defaultContext string, contexts []string) FormModel {
	m.editing = false
	m.editID = uuid.Nil
	m.err = ""
	m.focusIndex = fieldTitle
	m.contexts = contexts
	m.contextIndex = 0
	for i, c := range contexts {
		if c == defaultContext {
			m.contextIndex = i
			break
		}
	}

	m.title.SetValue("")
	m.category.SetValue("")
	m.comment.SetValue("")

	m.title.Focus()
	m.category.Blur()
	m.comment.Blur()

	return m
}

func (m FormModel) SetForEdit(todo domain.Todo, contexts []string) FormModel {
	m.editing = true
	m.editID = todo.ID
	m.err = ""
	m.focusIndex = fieldTitle
	m.contexts = contexts
	m.contextIndex = 0
	for i, c := range contexts {
		if c == todo.Context {
			m.contextIndex = i
			break
		}
	}

	m.title.SetValue(todo.Title)
	m.category.SetValue(todo.Category)
	m.comment.SetValue(todo.Comment)

	m.title.Focus()
	m.category.Blur()
	m.comment.Blur()

	return m
}

func (m FormModel) selectedContext() string {
	if len(m.contexts) == 0 {
		return ""
	}
	return m.contexts[m.contextIndex]
}

func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Save):
			return m, nil
		case key.Matches(msg, keys.Tab):
			m = m.nextField()
			return m, nil
		case key.Matches(msg, keys.ShiftTab):
			m = m.prevField()
			return m, nil
		}

		if m.focusIndex == fieldContext {
			switch {
			case key.Matches(msg, keys.Right), key.Matches(msg, keys.Down):
				if len(m.contexts) > 0 {
					m.contextIndex = (m.contextIndex + 1) % len(m.contexts)
				}
				return m, nil
			case key.Matches(msg, keys.Left), key.Matches(msg, keys.Up):
				if len(m.contexts) > 0 {
					m.contextIndex = (m.contextIndex - 1 + len(m.contexts)) % len(m.contexts)
				}
				return m, nil
			}
		}
	}

	return m.updateFocusedField(msg)
}

func (m FormModel) nextField() FormModel {
	m.focusIndex = (m.focusIndex + 1) % fieldCount
	m.applyFocus()
	return m
}

func (m FormModel) prevField() FormModel {
	m.focusIndex = (m.focusIndex - 1 + fieldCount) % fieldCount
	m.applyFocus()
	return m
}

func (m *FormModel) applyFocus() {
	m.title.Blur()
	m.category.Blur()
	m.comment.Blur()

	switch m.focusIndex {
	case fieldTitle:
		m.title.Focus()
	case fieldCategory:
		m.category.Focus()
	case fieldComment:
		m.comment.Focus()
	}
}

func (m FormModel) updateFocusedField(msg tea.Msg) (FormModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case fieldTitle:
		m.title, cmd = m.title.Update(msg)
	case fieldCategory:
		m.category, cmd = m.category.Update(msg)
	case fieldComment:
		m.comment, cmd = m.comment.Update(msg)
	}
	return m, cmd
}

func (m FormModel) Validate() string {
	if strings.TrimSpace(m.title.Value()) == "" {
		return "Title is required"
	}
	if m.selectedContext() == "" {
		return "Context is required — create a context first"
	}
	return ""
}

func (m FormModel) Values() (title, context, category, comment string) {
	return strings.TrimSpace(m.title.Value()),
		m.selectedContext(),
		strings.TrimSpace(m.category.Value()),
		m.comment.Value()
}

func (m FormModel) View() string {
	var b strings.Builder

	if m.editing {
		b.WriteString(formTitleStyle.Render("Edit Todo"))
	} else {
		b.WriteString(formTitleStyle.Render("New Todo"))
	}
	b.WriteString("\n\n")

	// Title
	titleIndicator := "  "
	if m.focusIndex == fieldTitle {
		titleIndicator = "▸ "
	}
	b.WriteString(titleIndicator + formLabelStyle.Render("Title") + "\n")
	b.WriteString("  " + m.title.View() + "\n\n")

	// Context selector
	ctxIndicator := "  "
	if m.focusIndex == fieldContext {
		ctxIndicator = "▸ "
	}
	b.WriteString(ctxIndicator + formLabelStyle.Render("Context") + "\n")
	b.WriteString("  " + m.renderContextSelector() + "\n\n")

	// Category
	catIndicator := "  "
	if m.focusIndex == fieldCategory {
		catIndicator = "▸ "
	}
	b.WriteString(catIndicator + formLabelStyle.Render("Category") + "\n")
	b.WriteString("  " + m.category.View() + "\n\n")

	// Comment
	cmtIndicator := "  "
	if m.focusIndex == fieldComment {
		cmtIndicator = "▸ "
	}
	b.WriteString(cmtIndicator + formLabelStyle.Render("Comment") + "\n")
	b.WriteString("  " + m.comment.View() + "\n\n")

	if m.err != "" {
		b.WriteString(errorStyle.Render(fmt.Sprintf("✗ %s", m.err)))
		b.WriteString("\n\n")
	}

	b.WriteString(dimStyle.Render("tab: next • shift+tab: prev • ctrl+s: save • esc: cancel"))

	return b.String()
}

func (m FormModel) renderContextSelector() string {
	if len(m.contexts) == 0 {
		return dimStyle.Render("(no contexts available)")
	}

	focused := m.focusIndex == fieldContext
	var parts []string
	for i, ctx := range m.contexts {
		if i == m.contextIndex {
			if focused {
				parts = append(parts, tabActive.Render(ctx))
			} else {
				parts = append(parts, detailTitleStyle.Render(ctx))
			}
		} else {
			parts = append(parts, dimStyle.Render(ctx))
		}
	}

	selector := strings.Join(parts, "  ")
	if focused {
		selector += "  " + dimStyle.Render("← →")
	}
	return selector
}
