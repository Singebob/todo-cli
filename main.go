package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"todo-cli/internal/app"
	"todo-cli/internal/infra/storage"
	"todo-cli/internal/ui"
)

func main() {
	store, err := storage.NewJSONStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	service := app.NewTodoService(store, store.ContextStore())
	model := ui.NewModel(service)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
