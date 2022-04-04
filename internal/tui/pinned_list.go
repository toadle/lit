package tui

import (
	"strings"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"launcher/internal/tui/list"
)

type PinnedListItemDelegate struct{}

func (d PinnedListItemDelegate) Height() int									{ return 1 }
func (d PinnedListItemDelegate) Spacing() int									{ return 0 }
func (d PinnedListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd	{ return nil }
func (d PinnedListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PinnedListItem)
	if !ok {
		return
	}
	fmt.Fprintf(w, i.AsFormattedString())
}

type PinnedListItem struct {
	result			string
	cmdStr			string
	itemFormat		string
	whenSelected	string
	currentInput	string
	successful		bool
	errorMsg		string
}

func (i PinnedListItem) AsFormattedString() string {
	var (
		sections    []string
	)
	if len(i.result) > 0 {
		sections = append(sections, styles.Text.Render(i.result))
		sections = append(sections, " ")
	}
	if i.successful {
		sections = append(sections, styles.SuccessText.Render("✔"))
	} else {
		sections = append(sections, styles.ErrorText.Render("✘"))
	}
	sections = append(sections, " ")
	cmdStr := i.cmdStr
	if len(i.currentInput) > 0 {
		cmdStr = strings.Replace(i.cmdStr, "{input}", styles.MutedTextUnterlined.Render(i.currentInput), 1)
	}
	sections = append(sections, styles.MutedText.Render(cmdStr))
	return styles.PinnedListItem.Render(lipgloss.JoinHorizontal(1, sections...))
}
func (i PinnedListItem) FilterValue() string { return i.cmdStr }
func (i *PinnedListItem) SetCurrentValue(str string) { i.currentInput = str }

func NewPinnedListItem(cmdStr, itemFormat, whenSelected string) PinnedListItem {
	return PinnedListItem{
		cmdStr: cmdStr,
		itemFormat: itemFormat,
		whenSelected: whenSelected,
		successful: false,
	}
}
