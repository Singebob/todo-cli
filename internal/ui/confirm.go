package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type confirmResult int

const (
	confirmNone confirmResult = iota
	confirmYes
	confirmNo
)

type ConfirmModel struct {
	message string
	todoID  uuid.UUID
	active  bool
	result  confirmResult
}

func NewConfirmModel() ConfirmModel {
	return ConfirmModel{}
}

func (m ConfirmModel) SetMessage(msg string, id uuid.UUID) ConfirmModel {
	m.message = msg
	m.todoID = id
	m.active = true
	m.result = confirmNone
	return m
}

func (m ConfirmModel) Update(msg tea.Msg) (ConfirmModel, tea.Cmd) {
	if !m.active {
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Yes):
			m.result = confirmYes
			m.active = false
		case key.Matches(msg, keys.No), key.Matches(msg, keys.Back):
			m.result = confirmNo
			m.active = false
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	if !m.active {
		return ""
	}
	content := m.message + "\n\n" + dimStyle.Render("[y]es  [n]o")
	return lipgloss.Place(0, 0,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(content),
	)
}
