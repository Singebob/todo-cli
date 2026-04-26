package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"todo-cli/internal/app"
	"todo-cli/internal/domain"
)

type view int

const (
	viewList view = iota
	viewDetail
	viewCreate
	viewEdit
	viewConfirm
	viewNewContext
)

type Model struct {
	service *app.TodoService

	activeView view

	tabBar    TabBarModel
	todoList  TodoListModel
	preview   PreviewModel
	detail    DetailModel
	form      FormModel
	confirm   ConfirmModel
	statusBar StatusBarModel

	contextInput textinput.Model

	currentContext string
	width, height  int
}

func NewModel(service *app.TodoService) Model {
	ci := textinput.New()
	ci.Placeholder = "Context name"
	ci.CharLimit = 40

	return Model{
		service:      service,
		activeView:   viewList,
		tabBar:       NewTabBarModel(),
		todoList:     NewTodoListModel(),
		preview:      NewPreviewModel(),
		detail:       NewDetailModel(),
		form:         NewFormModel(),
		confirm:      NewConfirmModel(),
		statusBar:    NewStatusBarModel(),
		contextInput: ci,
	}
}

// --- tea.Msg types ---

type contextsLoadedMsg struct{ contexts []string }
type todosLoadedMsg struct{ todos []domain.Todo }
type todoCreatedMsg struct{ todo domain.Todo }
type todoUpdatedMsg struct{ todo domain.Todo }
type todoDeletedMsg struct{ id uuid.UUID }
type todoToggledMsg struct{ todo domain.Todo }
type categoriesLoadedMsg struct{ categories []string }
type contextCreatedMsg struct{ name string }
type contextDeletedMsg struct{}
type errMsg struct{ err error }

// --- Commands ---

func (m Model) loadContexts() tea.Cmd {
	return func() tea.Msg {
		contexts, err := m.service.AllContexts()
		if err != nil {
			return errMsg{err}
		}
		return contextsLoadedMsg{contexts}
	}
}

func (m Model) loadTodos() tea.Cmd {
	return func() tea.Msg {
		if m.currentContext == "" {
			return todosLoadedMsg{nil}
		}
		var todos []domain.Todo
		var err error
		if m.todoList.showCompleted {
			todos, err = m.service.ListCompletedByContext(m.currentContext)
		} else if m.todoList.categoryFilter != "" {
			todos, err = m.service.ListByContextAndCategory(m.currentContext, m.todoList.categoryFilter)
		} else {
			todos, err = m.service.ListByContext(m.currentContext)
		}
		if err != nil {
			return errMsg{err}
		}
		return todosLoadedMsg{todos}
	}
}

func (m Model) createTodo(title, comment, context, category string) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.service.CreateTodo(title, comment, context, category)
		if err != nil {
			return errMsg{err}
		}
		return todoCreatedMsg{todo}
	}
}

func (m Model) updateTodo(id uuid.UUID, title, comment, context, category string) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.service.UpdateTodo(id, title, comment, context, category)
		if err != nil {
			return errMsg{err}
		}
		return todoUpdatedMsg{todo}
	}
}

func (m Model) deleteTodo(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteTodo(id); err != nil {
			return errMsg{err}
		}
		return todoDeletedMsg{id}
	}
}

func (m Model) toggleTodo(id uuid.UUID) tea.Cmd {
	return func() tea.Msg {
		todo, err := m.service.ToggleDone(id)
		if err != nil {
			return errMsg{err}
		}
		return todoToggledMsg{todo}
	}
}

func (m Model) createContext(name string) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.CreateContext(name); err != nil {
			return errMsg{err}
		}
		return contextCreatedMsg{name}
	}
}

func (m Model) deleteContext(name string) tea.Cmd {
	return func() tea.Msg {
		if err := m.service.DeleteContext(name); err != nil {
			return errMsg{err}
		}
		return contextDeletedMsg{}
	}
}

func (m Model) loadCategories() tea.Cmd {
	return func() tea.Msg {
		cats, err := m.service.CategoriesForContext(m.currentContext)
		if err != nil {
			return errMsg{err}
		}
		return categoriesLoadedMsg{cats}
	}
}

// --- Init ---

func (m Model) Init() tea.Cmd {
	return m.loadContexts()
}

// --- Update ---

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateLayout()
		return m, nil

	case contextsLoadedMsg:
		m.tabBar.contexts = msg.contexts
		if m.currentContext == "" && len(msg.contexts) > 0 {
			m.currentContext = msg.contexts[0]
			m.tabBar.cursor = 0
			return m, m.loadTodos()
		}
		if m.currentContext != "" {
			m.tabBar = m.tabBar.SelectByName(m.currentContext)
		}
		return m, nil

	case todosLoadedMsg:
		m.todoList.todos = msg.todos
		m.todoList.cursor = 0
		m.todoList.scrollOffset = 0
		m = m.syncPreview()
		return m, nil

	case todoCreatedMsg:
		m.currentContext = msg.todo.Context
		m.tabBar = m.tabBar.SelectByName(m.currentContext)
		m.activeView = viewList
		m = m.updateHints()
		return m, m.loadTodos()

	case todoUpdatedMsg:
		m.currentContext = msg.todo.Context
		m.tabBar = m.tabBar.SelectByName(m.currentContext)
		m.activeView = viewList
		m = m.updateHints()
		return m, tea.Batch(m.loadContexts(), m.loadTodos())

	case todoDeletedMsg:
		m.activeView = viewList
		m = m.updateHints()
		return m, m.loadTodos()

	case todoToggledMsg:
		if m.activeView == viewDetail {
			m.detail = m.detail.SetTodo(msg.todo, m.contentWidth(), m.contentHeight())
		}
		return m, m.loadTodos()

	case categoriesLoadedMsg:
		cats := msg.categories
		current := m.todoList.categoryFilter
		if len(cats) == 0 {
			m.todoList.categoryFilter = ""
		} else if current == "" {
			m.todoList.categoryFilter = cats[0]
		} else {
			m.todoList.categoryFilter = ""
			for i, c := range cats {
				if c == current && i+1 < len(cats) {
					m.todoList.categoryFilter = cats[i+1]
					break
				}
			}
		}
		return m, m.loadTodos()

	case contextCreatedMsg:
		m.currentContext = msg.name
		m.activeView = viewList
		m = m.updateHints()
		return m, tea.Batch(m.loadContexts(), m.loadTodos())

	case contextDeletedMsg:
		m.activeView = viewList
		m = m.updateHints()
		return m, m.loadContexts()

	case errMsg:
		m.statusBar.hints = errorStyle.Render("Error: " + msg.err.Error())
		return m, nil
	}

	switch m.activeView {
	case viewList:
		return m.updateList(msg)
	case viewDetail:
		return m.updateDetail(msg)
	case viewCreate, viewEdit:
		return m.updateForm(msg)
	case viewConfirm:
		return m.updateConfirm(msg)
	case viewNewContext:
		return m.updateNewContext(msg)
	}

	return m, nil
}

func (m Model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Tab):
			m.tabBar = m.tabBar.Next()
			m.currentContext = m.tabBar.SelectedContext()
			m.todoList.showCompleted = false
			m.todoList.categoryFilter = ""
			return m, m.loadTodos()

		case key.Matches(msg, keys.ShiftTab):
			m.tabBar = m.tabBar.Prev()
			m.currentContext = m.tabBar.SelectedContext()
			m.todoList.showCompleted = false
			m.todoList.categoryFilter = ""
			return m, m.loadTodos()

		case key.Matches(msg, keys.Enter):
			if todo, ok := m.todoList.SelectedTodo(); ok {
				m.activeView = viewDetail
				m.detail = m.detail.SetTodo(todo, m.contentWidth(), m.contentHeight())
				m = m.updateHints()
				return m, nil
			}

		case key.Matches(msg, keys.Space):
			if todo, ok := m.todoList.SelectedTodo(); ok {
				return m, m.toggleTodo(todo.ID)
			}

		case key.Matches(msg, keys.New):
			if m.currentContext == "" {
				m.activeView = viewNewContext
				m.contextInput.SetValue("")
				m.contextInput.Focus()
				m = m.updateHints()
				return m, nil
			}
			m.activeView = viewCreate
			m.form = m.form.SetForCreate(m.currentContext, m.tabBar.contexts)
			m = m.updateHints()
			return m, m.form.title.Focus()

		case key.Matches(msg, keys.Edit):
			if todo, ok := m.todoList.SelectedTodo(); ok {
				m.activeView = viewEdit
				m.form = m.form.SetForEdit(todo, m.tabBar.contexts)
				m = m.updateHints()
				return m, m.form.title.Focus()
			}

		case key.Matches(msg, keys.Delete):
			if todo, ok := m.todoList.SelectedTodo(); ok {
				m.activeView = viewConfirm
				m.confirm = m.confirm.SetMessage("Delete \""+todo.Title+"\"?", todo.ID)
				m = m.updateHints()
				return m, nil
			}

		case key.Matches(msg, keys.Completed):
			m.todoList.showCompleted = !m.todoList.showCompleted
			m.todoList.categoryFilter = ""
			return m, m.loadTodos()

		case key.Matches(msg, keys.Filter):
			return m, m.loadCategories()

		case key.Matches(msg, keys.NewContext):
			m.activeView = viewNewContext
			m.contextInput.SetValue("")
			m.contextInput.Focus()
			m = m.updateHints()
			return m, nil

		case key.Matches(msg, keys.DelContext):
			if m.currentContext != "" {
				return m, m.deleteContext(m.currentContext)
			}
		}
	}

	prevCursor := m.todoList.cursor
	var cmd tea.Cmd
	m.todoList, cmd = m.todoList.Update(msg)
	if m.todoList.cursor != prevCursor {
		m = m.syncPreview()
	}
	return m, cmd
}

func (m Model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.activeView = viewList
			m = m.updateHints()
			return m, nil

		case key.Matches(msg, keys.Edit):
			m.activeView = viewEdit
			m.form = m.form.SetForEdit(m.detail.todo, m.tabBar.contexts)
			m = m.updateHints()
			return m, m.form.title.Focus()

		case key.Matches(msg, keys.Space):
			return m, m.toggleTodo(m.detail.todo.ID)

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.activeView == viewEdit {
				m.activeView = viewDetail
			} else {
				m.activeView = viewList
			}
			m = m.updateHints()
			return m, nil

		case key.Matches(msg, keys.Save):
			if errStr := m.form.Validate(); errStr != "" {
				m.form.err = errStr
				return m, nil
			}
			title, ctx, cat, cmt := m.form.Values()
			if m.form.editing {
				return m, m.updateTodo(m.form.editID, title, cmt, ctx, cat)
			}
			return m, m.createTodo(title, cmt, ctx, cat)
		}
	}

	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.confirm, _ = m.confirm.Update(msg)
	switch m.confirm.result {
	case confirmYes:
		id := m.confirm.todoID
		m.activeView = viewList
		m = m.updateHints()
		return m, m.deleteTodo(id)
	case confirmNo:
		m.activeView = viewList
		m = m.updateHints()
	}
	return m, nil
}

func (m Model) updateNewContext(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			m.activeView = viewList
			m = m.updateHints()
			return m, nil
		case key.Matches(msg, keys.Enter):
			name := m.contextInput.Value()
			if name == "" {
				return m, nil
			}
			return m, m.createContext(name)
		}
	}

	var cmd tea.Cmd
	m.contextInput, cmd = m.contextInput.Update(msg)
	return m, cmd
}

// --- Layout helpers ---

func (m Model) listWidth() int {
	return m.width / 2
}

func (m Model) previewWidth() int {
	return m.width - m.listWidth()
}

func (m Model) updateLayout() Model {
	m.tabBar.width = m.width
	m.todoList.width = m.listWidth()
	m.todoList.height = m.contentHeight()
	m.todoList.focused = true
	m.statusBar.width = m.width
	m = m.syncPreview()
	m = m.updateHints()
	return m
}

func (m Model) syncPreview() Model {
	if todo, ok := m.todoList.SelectedTodo(); ok {
		m.preview = m.preview.SetTodo(&todo, m.previewWidth(), m.contentHeight())
	} else {
		m.preview = m.preview.SetTodo(nil, m.previewWidth(), m.contentHeight())
	}
	return m
}

func (m Model) contentWidth() int {
	return m.width - 4
}

func (m Model) contentHeight() int {
	return m.height - 4
}

func (m Model) updateHints() Model {
	switch m.activeView {
	case viewList:
		m.statusBar.hints = "n:new  enter:open  space:done  e:edit  d:delete  c:completed  /:filter  tab/shift+tab:context  C:new ctx  D:del ctx  q:quit"
	case viewDetail:
		m.statusBar.hints = "e:edit  space:toggle done  esc:back  q:quit"
	case viewCreate, viewEdit:
		m.statusBar.hints = "tab:next  shift+tab:prev  ctrl+s:save  esc:cancel"
	case viewConfirm:
		m.statusBar.hints = "y:yes  n:no  esc:cancel"
	case viewNewContext:
		m.statusBar.hints = "enter:create  esc:cancel"
	}
	return m
}

// --- View ---

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.activeView {
	case viewDetail:
		content := m.detail.View()
		return lipgloss.JoinVertical(lipgloss.Left,
			m.tabBar.View(),
			lipgloss.NewStyle().Height(m.contentHeight()).Render(content),
			m.statusBar.View(),
		)

	case viewCreate, viewEdit:
		content := m.form.View()
		return lipgloss.JoinVertical(lipgloss.Left,
			m.tabBar.View(),
			lipgloss.NewStyle().Height(m.contentHeight()).Padding(1, 2).Render(content),
			m.statusBar.View(),
		)

	case viewConfirm:
		overlay := m.confirm.View()
		return lipgloss.JoinVertical(lipgloss.Left,
			m.tabBar.View(),
			lipgloss.NewStyle().Height(m.contentHeight()).Render(
				lipgloss.Place(m.width, m.contentHeight(),
					lipgloss.Center, lipgloss.Center,
					overlay,
					lipgloss.WithWhitespaceChars(" "),
				),
			),
			m.statusBar.View(),
		)

	case viewNewContext:
		prompt := formTitleStyle.Render("New Context") + "\n\n" + m.contextInput.View()
		box := promptBoxStyle.Render(prompt)
		return lipgloss.JoinVertical(lipgloss.Left,
			m.tabBar.View(),
			lipgloss.NewStyle().Height(m.contentHeight()).Render(
				lipgloss.Place(m.width, m.contentHeight(),
					lipgloss.Center, lipgloss.Center,
					box,
					lipgloss.WithWhitespaceChars(" "),
				),
			),
			m.statusBar.View(),
		)
	}

	listAndPreview := lipgloss.JoinHorizontal(lipgloss.Top,
		m.todoList.View(),
		m.preview.View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left,
		m.tabBar.View(),
		listAndPreview,
		m.statusBar.View(),
	)
}
