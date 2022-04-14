package tui

import (
	// "fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/tui/list"
	"lit/internal/tui/textinput"
	"lit/internal/tui/style"
	"lit/internal/config"
	"lit/internal/shell"
)

type queryChangedMsg int

type Bubble struct {
	styles			style.Styles
	config			*config.LauncherConfig
	width			int
	height			int

	resultList		list.Model
	queryInput		textinput.Model
	pinnedList		list.Model

	queryInputTag	int
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		styles: 	style.DefaultStyles(),
		config:		cliCfg,
		resultList:	list.New([]list.Item{}, 0),
		queryInput:	textinput.New(),
		pinnedList:	list.New([]list.Item{}, 0),
	}

	b.resultList.SetNoResultText("Nothing found.")

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
		pinnedItems = append(pinnedItems, list.NewPinnedListItem(sourceConfig.Command, sourceConfig.ItemFormat, sourceConfig.WhenSelected))
	}
	b.pinnedList.SetItems(pinnedItems)

	return tea.Batch(teaCmds...)
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var teaCmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return b, tea.Quit
		case tea.KeyEnter:
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
		case tea.KeyRunes, tea.KeyBackspace:
			b.queryInputTag++
			teaCmds = append(teaCmds, tea.Tick(time.Millisecond * 100, func(_ time.Time) tea.Msg {
				return queryChangedMsg(b.queryInputTag)
			}))
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
	case queryChangedMsg:
		if int(msg) == b.queryInputTag {
			teaCmds = append(teaCmds, b.handleQueryChanged()...)
		}
	}

	var cmd tea.Cmd

	b.queryInput, cmd = b.queryInput.Update(msg)
	teaCmds = append(teaCmds, cmd)

	b.resultList, cmd = b.resultList.Update(msg)
	teaCmds = append(teaCmds, cmd)

	b.pinnedList, cmd = b.pinnedList.Update(msg)
	teaCmds = append(teaCmds, cmd)

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) handleQueryChanged() []tea.Cmd {
	var teaCmds []tea.Cmd

	currentInputValue := b.queryInput.Value()

	newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
		p:= i.(list.PinnedListItem)
		p.SetCurrentValue(currentInputValue)
		return p
	})
	b.pinnedList.SetItems(newPinnedList)

	for _, sourceConfig := range b.config.SourceConfigList {
		if strings.Contains(sourceConfig.Command, "{input}") {
			shellCmd := shell.NewCommand(sourceConfig.Command)
			shellCmd.SetInput(currentInputValue)
			teaCmds = append(teaCmds, shellCmd.Run)
		}
	}

	if (len(currentInputValue) == 0) {
		b.resultList.UnfilterItems()
		b.resultList.Unselect()
		b.queryInput.SetAutoComplete("")
	} else {
		b.resultList.SetFilterValue(currentInputValue)
		b.resultList.FilterItems()
		if len(b.resultList.VisibleItems()) > 0 {
			b.resultList.Select(0)
			b.queryInput.SetAutoComplete(b.resultList.VisibleItems()[0].FilterValue())
		}
	}
	return teaCmds
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
