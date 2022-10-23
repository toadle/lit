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

type CalculatorListItem struct {
	styles       style.Styles
	output       string
	currentInput string
	successful   bool

	sourceConfig config.CalculatorConfig
}

func (d CalculatorListItem) Update(msg tea.Msg, m *Model) tea.Cmd {
	return nil
}
func (d CalculatorListItem) Render(w io.Writer, m Model, index int, listItem Item) {
	i, ok := listItem.(CalculatorListItem)
	if !ok {
		return
	}
	var sections []string

	var textStyle, mutedTextStyle lipgloss.Style
	if index == m.Index() {
		textStyle = d.styles.CalculatorText
		mutedTextStyle = d.styles.MutedCalculatorText
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
	if idx := strings.Index(label, "{input}"); idx > -1 {
		label = strings.Replace(label, "{input}", currentInput, 1)
		underlineTextStyle := mutedTextStyle.Copy().Underline(true)
		label = lipgloss.StyleRunes(label, lo.RangeFrom(idx, idx+lengthOfInput), underlineTextStyle, mutedTextStyle)
	}
	sections = append(sections, mutedTextStyle.Render(label))
	sections = append(sections, " ")
	if len(i.output) > 0 {
		sections = append(sections, textStyle.Render(i.cleanedOutput()))
	}

	fmt.Fprintf(w, i.styles.CalculatorListItem.Render(lipgloss.JoinHorizontal(1, sections...)))
}
func (i CalculatorListItem) FilterValue() string         { return i.sourceConfig.Command }
func (i CalculatorListItem) CmdStr() string              { return i.sourceConfig.Command }
func (i CalculatorListItem) Action() string              { return i.sourceConfig.Action }
func (i *CalculatorListItem) SetCurrentValue(str string) { i.currentInput = str }
func (i *CalculatorListItem) SetOutput(str string)       { i.output = str }
func (i *CalculatorListItem) Output() string             { return i.output }
func (i *CalculatorListItem) SetSuccessful(b bool)       { i.successful = b }
func (i CalculatorListItem) cleanedOutput() string {
	return strings.Replace(i.output, "\n", "", -1)
}

func NewCalculatorListItem(sc config.CalculatorConfig) CalculatorListItem {
	return CalculatorListItem{
		styles:       style.DefaultStyles(),
		sourceConfig: sc,
		successful:   false,
	}
}
