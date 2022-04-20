package list

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/tui/style"
	"lit/internal/util"
)

type ResultListItem struct {
	styles			style.Styles
	resultData		util.ParsedResult
	whenSelected 	string
}

func (i ResultListItem) label() string {
	label, exists := i.resultData.Data["label"]
	if exists {
		return label
	} else {
		return ""
	}
}
func (i ResultListItem) data() string {
	data, exists := i.resultData.Data["data"]
	if exists {
		return data
	} else {
		return ""
	}
}

func (i ResultListItem) FilterValue() string {
	return i.label()
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

	if (m.filterState == Filtered) {
		underlineTextStyle := textStyle.Copy().Underline(true)
		matchedRunes := m.MatchesForItem(index)
		label := lipgloss.StyleRunes(i.label(), matchedRunes, underlineTextStyle, textStyle)
		sections = append(sections, textStyle.Render(label))
	} else {
		sections = append(sections, textStyle.Render(i.label()))
	}

	sections = append(sections, " ")
	sections = append(sections, mutedTextStyle.Render(i.data()))

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}

func NewResultListItem(itemData, itemFormat, whenSelected string) ResultListItem {
	parsedResult := util.ParseResult(itemData, itemFormat)
	return ResultListItem{
		styles: style.DefaultStyles(),
		resultData: parsedResult,
		whenSelected: whenSelected,
	}
}
