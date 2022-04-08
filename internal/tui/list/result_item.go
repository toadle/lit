package list

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/tui/style"
)

type ResultListItem struct {
	styles			style.Styles
	label			string
	data			string
	itemFormat		string
	whenSelected 	string
}

func (i ResultListItem) FilterValue() string {
	return i.label
}
func (d ResultListItem) Update(msg tea.Msg, m *Model) tea.Cmd	{ return nil }
func (d ResultListItem) Render(w io.Writer, m Model, index int, listItem Item) {
	i, ok := listItem.(ResultListItem)
	if !ok {
		return
	}

	var sections []string
	var textStyle, mutedTextStyle lipgloss.Style
	if index == m.Index() {
		textStyle = d.styles.SelectedText
		mutedTextStyle = d.styles.SelectedMutedText
	} else {
		textStyle = d.styles.Text
		mutedTextStyle = d.styles.MutedText
	}

	sections = append(sections, textStyle.Render(i.label))
	sections = append(sections, " ")
	sections = append(sections, mutedTextStyle.Render(i.data))

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}

func NewResultListItem(itemData, itemFormat, whenSelected string) ResultListItem {
	var label, data string
	parts := strings.Split(itemData, ":")
	if len(parts) == 2 {
		data = parts[0]
		label = parts[1]
	}
	return ResultListItem{
		styles: style.DefaultStyles(),
		label: label,
		data: data,
		itemFormat: itemFormat,
		whenSelected: whenSelected,
	}
}
