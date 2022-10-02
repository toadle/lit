package tui

import (
	// "fmt"

	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/config"
	"lit/internal/shell"
	"lit/internal/tui/list"
	"lit/internal/tui/style"
)

type queryChangedMsg int

type CursorFocus int

const (
	PinnedList CursorFocus = iota
	QueryInput
	ResultList
)

type Bubble struct {
	styles style.Styles
	config *config.LauncherConfig
	width  int
	height int

	resultList list.Model
	queryInput textinput.Model
	pinnedList list.Model

	queryInputTag int
	focus         CursorFocus
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		styles:     style.DefaultStyles(),
		config:     cliCfg,
		resultList: list.New([]list.Item{}, 0),
		queryInput: textinput.New(),
		pinnedList: list.New([]list.Item{}, 0),
		focus:      QueryInput,
	}

	b.resultList.SetNoResultText("Nothing found.")
	b.queryInput.Placeholder = "Your Query"
	b.queryInput.ShowCompletions = true
	b.focusQueryInput()

	return b
}

func (b *Bubble) Init() tea.Cmd {
	var teaCmds []tea.Cmd

	for _, sourceConfig := range b.config.MultiLineConfigList {
		shellCmd := shell.NewCommand(sourceConfig.Command, true)
		teaCmds = append(teaCmds, shellCmd.Run)
	}

	var pinnedItems []list.Item
	for _, sourceConfig := range b.config.SingleLineConfigList {
		pinnedItems = append(pinnedItems, list.NewPinnedListItem(sourceConfig))
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
		case tea.KeyUp:
			switch b.focus {
			case QueryInput:
				b.focusPinnedList()
			case ResultList:
				if b.resultList.Index() == 0 {
					b.focusQueryInput()
				}
			}
		case tea.KeyDown:
			switch b.focus {
			case PinnedList:
				if b.pinnedList.Index() == b.pinnedList.Height() {
					b.focusQueryInput()
				}
			case QueryInput:
				b.focusResultList()
			}
		case tea.KeyRunes, tea.KeyBackspace:
			b.queryInputTag++
			teaCmds = append(teaCmds, tea.Tick(time.Millisecond*100, func(_ time.Time) tea.Msg {
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
		sourceConfig, ok := b.config.SourceConfigFor(msg.CmdStr, msg.Multiline)
		if ok {
			if msg.Multiline {
				items := b.resultList.Items()
				for _, line := range msg.Lines() {
					items = append(items, list.NewResultListItem(line, sourceConfig))
				}
				b.resultList.SetItems(items)
			} else {
				newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
					p := i.(list.PinnedListItem)
					if p.CmdStr() == msg.CmdStr {
						p.SetOutput(msg.Output)
						p.SetSuccessful(msg.Successful)
					}
					return p
				})
				b.pinnedList.SetItems(newPinnedList)
			}
		}
	case queryChangedMsg:
		if int(msg) == b.queryInputTag {
			teaCmds = append(teaCmds, b.handleQueryChanged()...)
		}
	}

	var cmd tea.Cmd

	if b.focus == QueryInput {
		b.queryInput, cmd = b.queryInput.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == ResultList {
		b.resultList, cmd = b.resultList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == PinnedList {
		b.pinnedList, cmd = b.pinnedList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) handleQueryChanged() []tea.Cmd {
	var teaCmds []tea.Cmd

	currentInputValue := b.queryInput.Value()

	newPinnedList := lo.Map[list.Item, list.Item](b.pinnedList.Items(), func(i list.Item, _ int) list.Item {
		p := i.(list.PinnedListItem)
		p.SetCurrentValue(currentInputValue)
		return p
	})
	b.pinnedList.SetItems(newPinnedList)

	commands := b.generateCommandsFor(b.config.MultiLineConfigList, true)
	commands = append(commands, b.generateCommandsFor(b.config.SingleLineConfigList, false)...)

	commandsAsTeaCmds := lo.Map[shell.Command, tea.Cmd](commands, func(c shell.Command, _ int) tea.Cmd {
		return c.Run
	})

	teaCmds = append(teaCmds, commandsAsTeaCmds...)

	if len(currentInputValue) == 0 {
		b.resultList.UnfilterItems()
		b.resultList.Unselect()
	} else {
		b.resultList.SetFilterValue(currentInputValue)
		b.resultList.FilterItems()
		if len(b.resultList.VisibleItems()) > 0 {
			b.resultList.Select(0)
		}
	}

	teaCmds = append(teaCmds, b.generateCompletions)
	return teaCmds
}

func (b Bubble) generateCommandsFor(configs []config.SourceConfig, multiline bool) []shell.Command {
	currentInputValue := b.queryInput.Value()
	shellCommands := []shell.Command{}
	for _, sourceConfig := range configs {
		if strings.Contains(sourceConfig.Command, "{input}") {
			shellCmd := shell.NewCommand(sourceConfig.Command, multiline)
			shellCmd.SetParams(map[string]string{
				"input": currentInputValue,
			})
			shellCommands = append(shellCommands, *shellCmd)
		}
	}
	return shellCommands
}

func (b *Bubble) generateCompletions() tea.Msg {
	if len(b.queryInput.Value()) == 0 {
		return b.queryInput.NewCompletionMsg("")
	} else {
		if len(b.resultList.VisibleItems()) > 0 {
			autocompleteValue := b.resultList.VisibleItems()[0].FilterValue()
			return b.queryInput.NewCompletionMsg(autocompleteValue)
		}
	}
	return nil
}

func (b *Bubble) View() string {
	var (
		sections []string
	)

	queryStyle := b.styles.Query.Width(b.width)

	sections = append(sections, b.pinnedList.View())
	sections = append(sections, queryStyle.Render(b.queryInput.View()))
	sections = append(sections, b.resultList.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (b *Bubble) focusPinnedList() {
	b.focus = PinnedList
	b.resultList.Unselect()
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text
	b.pinnedList.Select(b.pinnedList.Height())
}
func (b *Bubble) focusQueryInput() {
	b.focus = QueryInput
	b.queryInput.Focus()
	b.queryInput.PromptStyle = b.styles.QueryPromptFocused
	b.pinnedList.Unselect()
	b.resultList.Unselect()
}
func (b *Bubble) focusResultList() {
	b.focus = ResultList
	b.pinnedList.Unselect()
	b.resultList.Select(0)
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text
}
