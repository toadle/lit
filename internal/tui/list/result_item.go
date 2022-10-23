package list

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/config"
	"lit/internal/shell"
	"lit/internal/tui/style"
	"lit/internal/util"
)

type ResultListItem struct {
	styles       style.Styles
	resultData   shell.CommandResult
	sourceConfig config.MultiLineSourceConfig
}

func (i ResultListItem) title() string {
	labelFormatStr := "{title}"
	labelsConfig := i.sourceConfig.Labels

	if len(labelsConfig.Title) > 0 {
		labelFormatStr = labelsConfig.Title
	}
	return shell.SetCommandParameters(labelFormatStr, i.resultData.Params)
}

func (i ResultListItem) description() string {
	labelFormatStr := "{description}"
	labelsConfig := i.sourceConfig.Labels

	if len(labelsConfig.Description) > 0 {
		labelFormatStr = labelsConfig.Description
	}
	return shell.SetCommandParameters(labelFormatStr, i.resultData.Params)
}

func (i ResultListItem) Action() string {
	return i.sourceConfig.Action
}

func (i ResultListItem) Params() map[string]string {
	return i.resultData.Params
}

func (i ResultListItem) FilterValue() string {
	return util.RemoveSpecialCharacters(i.title())
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
		matchedRunes := util.FindAllOccurrencesOfCharacters(i.title(), m.filterValue)
		label := lipgloss.StyleRunes(i.title(), matchedRunes, underlineTextStyle, textStyle)
		sections = append(sections, textStyle.Render(label))
	} else {
		sections = append(sections, textStyle.Render(i.title()))
	}

	sections = append(sections, " ")
	sections = append(sections, mutedTextStyle.Render(i.description()))

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}

func NewResultListItem(itemData string, sourceConfig config.MultiLineSourceConfig) ResultListItem {
	parsedResult := shell.ParseCommandResult(itemData, sourceConfig.Format)
	return ResultListItem{
		styles:       style.DefaultStyles(),
		resultData:   parsedResult,
		sourceConfig: sourceConfig,
	}
}
