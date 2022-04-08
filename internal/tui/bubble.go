package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/tui/list"
	"lit/internal/tui/style"
	"lit/internal/config"
	"lit/internal/shell"
)
type Bubble struct {
	styles		style.Styles
	config		*config.LauncherConfig
	width		int
	height		int

	resultList	list.Model
	queryInput	textinput.Model
	pinnedList	list.Model
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		styles: 	style.DefaultStyles(),
		config:		cliCfg,
		resultList:	list.New([]list.Item{}, 0),
		queryInput:	textinput.New(),
		pinnedList:	list.New([]list.Item{}, 0),
	}
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

	b.resultList.Select(0)

	var pinnedItems []list.Item
	for _, sourceConfig := range b.config.PinnedSourceConfigList() {
		pinnedItems = append(pinnedItems, list.NewPinnedListItem(sourceConfig.Command, sourceConfig.ItemFormat, sourceConfig.WhenSelected))
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
// 			i, ok := b.resultList.SelectedItem().(*ResultListItem)
// 			if ok {
// 				cmdStr := strings.Replace(i.whenSelected, "{data}", i.data, 1)
// 				commandComponents := strings.Split(strings.TrimSpace(cmdStr)," ")
// 				mainCommand := commandComponents[0]
// 				args := commandComponents[1:]
//
// 				cmd := exec.Command(mainCommand, args...)
// 				err := cmd.Run()
//
// 				if err != nil {
// 					fmt.Println(err)
// 				} else {
// 					os.Exit(0)
// 				}
// 			}
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		// _, right, _, left := styles.App.GetMargin()

		b.resultList.SetHeight(7)
		b.pinnedList.SetHeight(len(b.pinnedList.Items()))
	case shell.ShellCommandResultMsg:
		sourceConfig, ok := b.config.SourceConfigFor(msg.CmdStr)
		if ok {
			if sourceConfig.Pinned {
				newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
					p:= i.(list.PinnedListItem)
					if p.CmdStr() == msg.CmdStr {
						p.SetOutput(msg.Output)
						p.SetSuccessful(msg.Successful)
					}
					return p
				})
				b.pinnedList.SetItems(newPinnedList)
			} else {
				items := b.resultList.Items()
				for _, line := range msg.Lines() {
					items = append(items, list.NewResultListItem(line, sourceConfig.ItemFormat, sourceConfig.WhenSelected))
				}
				b.resultList.SetItems(items)
			}
		}
	}

	// Update the filter text input component
	newQueryInput, cmd := b.queryInput.Update(msg)
	queryChanged := b.queryInput.Value() != newQueryInput.Value()
	b.queryInput = newQueryInput
	teaCmds = append(teaCmds, cmd)

	if queryChanged {
		newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
			p:= i.(list.PinnedListItem)
			p.SetCurrentValue(newQueryInput.Value())
			return p
		})
		b.pinnedList.SetItems(newPinnedList)

		for _, sourceConfig := range b.config.SourceConfigList {
			if strings.Contains(sourceConfig.Command, "{input}") {
				shellCmd := shell.NewCommand(sourceConfig.Command)
				shellCmd.SetInput(newQueryInput.Value())
				teaCmds = append(teaCmds, shellCmd.Run)
			}
		}
	}

	list, cmd := b.resultList.Update(msg)
	b.resultList = list
	teaCmds = append(teaCmds, cmd)

	list, cmd = b.pinnedList.Update(msg)
	b.pinnedList = list
	teaCmds = append(teaCmds, cmd)

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) View() string {
	var (
		sections    []string
	)

	queryStyle := b.styles.Query.Width(b.width)

	sections = append(sections, b.pinnedList.View())
	sections = append(sections, queryStyle.Render(b.queryInput.View()))
	sections = append(sections, b.resultList.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
