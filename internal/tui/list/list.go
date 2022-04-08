package list

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/tui/style"
)


type Model struct {
	height				int
	items				[]Item
	styles				style.Styles

	cursor				int
	windowBeginIndex	int
	windowEndIndex		int
}

func New(items []Item, height int) Model {
	m := Model{
		cursor: 			-1,
		styles: 			style.DefaultStyles(),
		height:				height,
		windowBeginIndex: 	0,
		windowEndIndex: 	0,
		items:				items,
	}
	return m
}

func (m *Model) Select(index int) {
	m.cursor = index
}

func (m *Model) Unselect() {
	m.cursor = -1
}

func (m *Model) Index() int {
	return m.cursor
}

func (m *Model) SetItems(i []Item) {
	m.items = i
}
func (m Model) Items() []Item {
	return m.items
}

func (m *Model) SetHeight(h int) {
	m.height = h
	if m.windowEndIndex < m.height {
		m.windowEndIndex = h
	}
}

func (m *Model) CursorUp() {
	if m.cursor > -1 {
		if m.cursor != 0 {
			m.cursor--
		}

		if m.cursor < m.windowBeginIndex {
			m.windowBeginIndex--
			m.windowEndIndex--
		}
	}
}

func (m *Model) CursorDown() {
	if m.cursor > -1 {
		if m.cursor == m.windowEndIndex - 1 && m.windowEndIndex < len(m.items) - 1 {
			m.windowBeginIndex++
			m.windowEndIndex++
		}
		m.cursor++
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var teaCmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.CursorUp()
		case "down":
			m.CursorDown()
		}
	}

	return m, tea.Batch(teaCmds...)

}

func (m Model) View() string {
	var sections    []string
	sections = append(sections, m.populatedView())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) visibleItems() []Item {
	return m.items
}

func (m Model) populatedView() string {
	items := m.visibleItems()

	var b strings.Builder

	// Empty states
	if len(items) == 0 {
		return m.styles.MutedText.Render("No items found.")
	}

	if len(items) > 0 {
		for i, item := range items[m.windowBeginIndex:m.windowEndIndex] {
			item.Render(&b, m, m.windowBeginIndex + i, item)
			fmt.Fprint(&b, "\n")
		}
	}

	return b.String()
}
