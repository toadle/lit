package list

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/config"
	"lit/internal/shell"
	"lit/internal/tui/style"
)

type ResultListItem struct {
	styles       style.Styles
	resultData   shell.CommandResult
	sourceConfig config.SourceConfig
}

func (i ResultListItem) label() string {
	label, exists := i.resultData.Params["label"]
	if exists {
		return label
	} else {
		return ""
	}
}
func (i ResultListItem) data() string {
	data, exists := i.resultData.Params["data"]
	if exists {
		return data
	} else {
		return ""
	}
}

func (i ResultListItem) Action() string {
	return i.sourceConfig.Action
}

func (i ResultListItem) Params() map[string]string {
	return i.resultData.Params
}

func (i ResultListItem) FilterValue() string {
	return i.label()
}
func (d ResultListItem) Update(msg tea.Msg, m *Model) tea.Cmd { return nil }
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

	if m.filterState == Filtered {
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

func NewResultListItem(itemData string, sourceConfig config.SourceConfig) ResultListItem {
	parsedResult := shell.ParseCommandResult(itemData, sourceConfig.Format)
	return ResultListItem{
		styles:       style.DefaultStyles(),
		resultData:   parsedResult,
		sourceConfig: sourceConfig,
	}
}
