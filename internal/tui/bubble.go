package tui

import (
	"strings"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"launcher/internal/tui/list"
	"launcher/internal/tui/style"
	"launcher/internal/config"
	"launcher/internal/shell"
)

var (
	styles = style.DefaultStyles()
)

type Bubble struct {
	config		*config.LauncherConfig
	width		int
	height		int

	resultList	list.Model
	queryInput	textinput.Model
	pinnedList	list.Model
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		config:		cliCfg,
		resultList:	list.New([]list.Item{}, ResultListItemDelegate{}, 0, 0),
		queryInput:	textinput.New(),
		pinnedList:	list.New([]list.Item{}, PinnedListItemDelegate{}, 0, 0),
	}
	b.resultList.SetShowHelp(false)
	b.resultList.SetShowStatusBar(false)
	b.resultList.SetShowTitle(false)
	b.resultList.SetShowFilter(false)

	b.pinnedList.SetShowHelp(false)
	b.pinnedList.SetShowStatusBar(false)
	b.pinnedList.SetShowTitle(false)
	b.pinnedList.SetShowPagination(false)
	b.pinnedList.SetShowFilter(false)

	b.queryInput.Placeholder = "Your Query"
	b.queryInput.Focus()

	return b
}

func (b *Bubble) Init() tea.Cmd {
	var teaCmds []tea.Cmd

	for _, sourceConfig := range b.config.ResultSourceConfigList() {
		shellCmd := shell.NewCommand(sourceConfig.Command)
		teaCmds = append(teaCmds, shellCmd.Run)
	}

	var pinnedItems []list.Item
	for _, sourceConfig := range b.config.PinnedSourceConfigList() {
		pinnedItems = append(pinnedItems, NewPinnedListItem(sourceConfig.Command, sourceConfig.ItemFormat, sourceConfig.WhenSelected))
	}
	b.pinnedList.SetItems(pinnedItems)

	return tea.Batch(teaCmds...)
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teaCmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "enter", " ":
			i, ok := b.resultList.SelectedItem().(*ResultListItem)
			if ok {
				cmdStr := strings.Replace(i.whenSelected, "{data}", i.data, 1)
				commandComponents := strings.Split(strings.TrimSpace(cmdStr)," ")
				mainCommand := commandComponents[0]
				args := commandComponents[1:]

				cmd := exec.Command(mainCommand, args...)
				err := cmd.Run()

				if err != nil {
					fmt.Println(err)
				} else {
					os.Exit(0)
				}
			}
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		_, right, _, left := styles.App.GetMargin()

		b.resultList.SetSize(msg.Width-left-right, 7)
		b.pinnedList.SetSize(msg.Width-left-right, len(b.pinnedList.Items()))
	case shell.ShellCommandResultMsg:
		sourceConfig, ok := b.config.SourceConfigFor(msg.OriginalCommandStr)
		if ok {
			items := b.resultList.Items()
			for _, line := range msg.Lines {
				items = append(items, NewResultListItem(line, sourceConfig.ItemFormat, sourceConfig.WhenSelected))
			}
			b.resultList.SetItems(items)
		}
	}

	// Update the filter text input component
	newQueryInput, cmd := b.queryInput.Update(msg)
	queryChanged := b.queryInput.Value() != newQueryInput.Value()
	b.queryInput = newQueryInput
	teaCmds = append(teaCmds, cmd)

	if queryChanged {
		newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
			p:= i.(PinnedListItem)
			p.SetCurrentValue(newQueryInput.Value())
			return p
		})
		b.pinnedList.SetItems(newPinnedList)
	}

	// list, cmd := b.resultList.Update(msg)
	// b.resultList = list

	list, cmd := b.pinnedList.Update(msg)
	b.pinnedList = list
	teaCmds = append(teaCmds, cmd)

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) View() string {
	var (
		sections    []string
	)

	queryStyle := styles.Query.Width(b.width)

	sections = append(sections, b.pinnedList.View())
	sections = append(sections, queryStyle.Render(b.queryInput.View()))
	sections = append(sections, b.resultList.View())

	return styles.App.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
