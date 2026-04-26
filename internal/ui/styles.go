package ui

import "github.com/charmbracelet/lipgloss"

var (
	listPaneStyle = lipgloss.NewStyle().
			Padding(0, 2)

	todoNormal = lipgloss.NewStyle()

	todoSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	todoDone = lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.Color("240"))

	categoryTag = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("172")).
			Bold(true).
			Padding(0, 1)

	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")).
				MarginBottom(1)

	detailMetaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	formLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Width(12)

	formTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1)

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	dialogBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 3).
			Width(50)

	promptBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 3).
			Width(50)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
