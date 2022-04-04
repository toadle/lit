package tui

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"launcher/internal/tui/list"
)

type ResultListItemDelegate struct{}

func (d ResultListItemDelegate) Height() int	{ return 1 }
func (d ResultListItemDelegate) Spacing() int	{ return 0 }
func (d ResultListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd	{ return nil }
func (d ResultListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ResultListItem)
	if !ok {
		return
	}
	fmt.Fprintf(w, i.AsFormattedString())
}


type ResultListItem struct {
	label, data, itemFormat, whenSelected string
}
func (i ResultListItem) AsFormattedString() string {
	var (
		sections	[]string
	)
	sections = append(sections, styles.Text.Render(i.label))
	sections = append(sections, " ")
	sections = append(sections, styles.MutedText.Render(i.data))

	return styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...))
}
func (i ResultListItem) FilterValue() string {
	return i.label
}


func NewResultListItem(itemData, itemFormat, whenSelected string) ResultListItem {
	var label, data string
	parts := strings.Split(itemData, ":")
	if len(parts) == 2 {
		data = parts[0]
		label = parts[1]
	}
	return ResultListItem{label: label, data: data, itemFormat: itemFormat, whenSelected: whenSelected}
}
