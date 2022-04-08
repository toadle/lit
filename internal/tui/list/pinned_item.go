package list

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lit/internal/tui/style"
)

type PinnedListItem struct {
	styles			style.Styles
	output			string
	cmdStr			string
	itemFormat		string
	whenSelected	string
	currentInput	string
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

	if len(i.output) > 0 {
		sections = append(sections, i.styles.Text.Render(i.cleanedOutput()))
		sections = append(sections, " ")
	}
	if i.successful {
		sections = append(sections, i.styles.SuccessText.Render("✔"))
	} else {
		sections = append(sections, i.styles.ErrorText.Render("✘"))
	}
	sections = append(sections, " ")
	cmdStr := i.cmdStr
	if len(i.currentInput) > 0 {
		cmdStr = strings.Replace(i.cmdStr, "{input}", i.styles.MutedTextUnterlined.Render(i.currentInput), 1)
	}
	sections = append(sections, i.styles.MutedText.Render(cmdStr))

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

func NewPinnedListItem(cmdStr, itemFormat, whenSelected string) PinnedListItem {
	return PinnedListItem{
		styles: style.DefaultStyles(),
		cmdStr: cmdStr,
		itemFormat: itemFormat,
		whenSelected: whenSelected,
		successful: false,
	}
}
