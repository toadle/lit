package list

import (
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/config"
	"lit/internal/tui/style"
)

type SingleListItem struct {
	styles       style.Styles
	output       string
	currentInput string
	successful   bool

	sourceConfig config.SingleLineSourceConfig
}

func (d SingleListItem) Update(msg tea.Msg, m *Model) tea.Cmd {
	return nil
}
func (d SingleListItem) Render(w io.Writer, m Model, index int, listItem Item) {
	i, ok := listItem.(SingleListItem)
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
			label = lipgloss.StyleRunes(label, lo.RangeFrom(idx, idx+lengthOfInput), underlineTextStyle, mutedTextStyle)
		}
	}
	sections = append(sections, mutedTextStyle.Render(label))
	sections = append(sections, " ")
	if len(i.output) > 0 {
		sections = append(sections, textStyle.Render(i.cleanedOutput()))
	}

	fmt.Fprintf(w, i.styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}
func (i SingleListItem) FilterValue() string         { return i.sourceConfig.Command }
func (i SingleListItem) CmdStr() string              { return i.sourceConfig.Command }
func (i SingleListItem) Action() string              { return i.sourceConfig.Action }
func (i *SingleListItem) SetCurrentValue(str string) { i.currentInput = str }
func (i *SingleListItem) SetOutput(str string)       { i.output = str }
func (i *SingleListItem) Output() string             { return i.output }
func (i *SingleListItem) SetSuccessful(b bool)       { i.successful = b }
func (i SingleListItem) cleanedOutput() string {
	return strings.Replace(i.output, "\n", "", -1)
}

func NewSingleListItem(sc config.SingleLineSourceConfig) SingleListItem {
	return SingleListItem{
		styles:       style.DefaultStyles(),
		sourceConfig: sc,
		successful:   false,
	}
}
