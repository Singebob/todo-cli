package ui

type StatusBarModel struct {
	hints string
	width int
}

func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{}
}

func (m StatusBarModel) View() string {
	bar := statusBarStyle.Width(m.width).Render(m.hints)
	return bar
}
