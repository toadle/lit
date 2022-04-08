package list

import (
	"io"
	tea "github.com/charmbracelet/bubbletea"
)

type Item interface {
	FilterValue() string
	Render(w io.Writer, m Model, index int, item Item)
	Update(msg tea.Msg, m *Model) tea.Cmd
}
