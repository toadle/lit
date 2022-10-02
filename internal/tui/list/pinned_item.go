package list

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/config"
	"lit/internal/tui/style"
)

type PinnedListItem struct {
	styles			style.Styles
	output			string
	currentInput	string
	successful		bool

	sourceConfig	config.SourceConfig
}

func (d PinnedListItem) Update(msg tea.Msg, m *Model) tea.Cmd	{
	return nil
}
func (d PinnedListItem) Render(w io.Writer, m Model, index int, listItem Item) {
	i, ok := listItem.(PinnedListItem)
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

	sourceConfig := i.sourceConfig

	label := sourceConfig.Label

	if len(label) == 0 {
		label = sourceConfig.Command
	}
	currentInput := i.currentInput
	lengthOfInput := len(currentInput)
	if lengthOfInput > 0 {
		if idx := strings.Index(label, "{input}"); idx > -1 {
			label = strings.Replace(label, "{input}", currentInput, 1)
			underlineTextStyle := mutedTextStyle.Copy().Underline(true)
			label = lipgloss.StyleRunes(label, lo.RangeFrom(idx, idx + lengthOfInput), underlineTextStyle, mutedTextStyle)
		}
	}
	sections = append(sections, mutedTextStyle.Render(label))
	sections = append(sections, " ")
	if len(i.output) > 0 {
		sections = append(sections, textStyle.Render(i.cleanedOutput()))
	}

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}
func (i PinnedListItem) FilterValue() string { return i.sourceConfig.Command }
func (i PinnedListItem) CmdStr() string { return i.sourceConfig.Command }
func (i *PinnedListItem) SetCurrentValue(str string) { i.currentInput = str }
func (i *PinnedListItem) SetOutput(str string) { i.output = str }
func (i *PinnedListItem) SetSuccessful(b bool) { i.successful = b }
func (i PinnedListItem) cleanedOutput() string {
	return strings.Replace(i.output, "\n", "", -1)
}

func NewPinnedListItem(sc config.SourceConfig) PinnedListItem {
	return PinnedListItem{
		styles: style.DefaultStyles(),
		sourceConfig: sc,
		successful: false,
	}
}
