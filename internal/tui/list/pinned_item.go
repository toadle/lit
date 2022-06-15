package list

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/tui/style"
)

type PinnedListItem struct {
	styles			style.Styles
	output			string
	cmdStr			string
	itemFormat		string
	whenSelected	string
	currentInput	string
	label			string
	successful		bool
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

	displayStr := i.cmdStr
	if len(i.label) > 0 {
		displayStr = i.label
	}
	currentInput := i.currentInput
	lengthOfInput := len(currentInput)
	if lengthOfInput > 0 {
		if idx := strings.Index(displayStr, "{input}"); idx > -1 {
			displayStr = strings.Replace(displayStr, "{input}", currentInput, 1)
			underlineTextStyle := mutedTextStyle.Copy().Underline(true)
			displayStr = lipgloss.StyleRunes(displayStr, lo.RangeFrom(idx, idx + lengthOfInput), underlineTextStyle, mutedTextStyle)
		}
	}
	sections = append(sections, mutedTextStyle.Render(displayStr))
	sections = append(sections, " ")
	if len(i.output) > 0 {
		sections = append(sections, textStyle.Render(i.cleanedOutput()))
	}

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}
func (i PinnedListItem) FilterValue() string { return i.cmdStr }
func (i PinnedListItem) CmdStr() string { return i.cmdStr }
func (i *PinnedListItem) SetCurrentValue(str string) { i.currentInput = str }
func (i *PinnedListItem) SetOutput(str string) { i.output = str }
func (i *PinnedListItem) SetSuccessful(b bool) { i.successful = b }
func (i PinnedListItem) cleanedOutput() string {
	return strings.Replace(i.output, "\n", "", -1)
}

func NewPinnedListItem(cmdStr string, itemFormat string, label string, whenSelected string) PinnedListItem {
	return PinnedListItem{
		styles: style.DefaultStyles(),
		cmdStr: cmdStr,
		itemFormat: itemFormat,
		whenSelected: whenSelected,
		label: label,
		successful: false,
	}
}
