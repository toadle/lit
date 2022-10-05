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

type FocusChangeMsg struct {
	newFocus CursorFocus
}

const (
	SingleList CursorFocus = iota
	QueryInput
	MultiList
)

type Bubble struct {
	styles style.Styles
	config *config.LauncherConfig
	width  int
	height int

	multiList  list.Model
	queryInput textinput.Model
	singleList list.Model

	queryInputTag int
	focus         CursorFocus
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		styles:     style.DefaultStyles(),
		config:     cliCfg,
		multiList:  list.New([]list.Item{}, 0),
		queryInput: textinput.New(),
		singleList: list.New([]list.Item{}, 0),
		focus:      QueryInput,
	}

	b.multiList.SetNoResultText("Nothing found.")
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
	b.singleList.SetItems(pinnedItems)

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
				teaCmds = append(teaCmds, b.focusSingleList)
			case MultiList:
				if b.multiList.Index() == 0 {
					teaCmds = append(teaCmds, b.focusQueryInput)
				}
			}
		case tea.KeyDown:
			switch b.focus {
			case SingleList:
				if b.singleList.Index() == b.singleList.Height()-1 {
					teaCmds = append(teaCmds, b.focusQueryInput)
				}
			case QueryInput:
				teaCmds = append(teaCmds, b.focusMultiList)
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

		b.multiList.SetHeight(7)
		b.singleList.SetHeight(len(b.singleList.Items()))
	case shell.ShellCommandResultMsg:
		sourceConfig, ok := b.config.SourceConfigFor(msg.CmdStr, msg.Multiline)
		if ok {
			if msg.Multiline {
				items := b.multiList.Items()
				for _, line := range msg.Lines() {
					items = append(items, list.NewResultListItem(line, sourceConfig))
				}
				b.multiList.SetItems(items)
			} else {
				newPinnedList := lo.Map[list.Item, list.Item](b.singleList.Items(), func(i list.Item, _ int) list.Item {
					p := i.(list.PinnedListItem)
					if p.CmdStr() == msg.CmdStr {
						p.SetOutput(msg.Output)
						p.SetSuccessful(msg.Successful)
					}
					return p
				})
				b.singleList.SetItems(newPinnedList)
			}
		}
	case queryChangedMsg:
		if int(msg) == b.queryInputTag {
			teaCmds = append(teaCmds, b.handleQueryChanged()...)
		}
	case FocusChangeMsg:
		b.focus = FocusChangeMsg(msg).newFocus
	}

	var cmd tea.Cmd

	if b.focus == QueryInput {
		b.queryInput, cmd = b.queryInput.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == MultiList {
		b.multiList, cmd = b.multiList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == SingleList {
		b.singleList, cmd = b.singleList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) handleQueryChanged() []tea.Cmd {
	var teaCmds []tea.Cmd

	currentInputValue := b.queryInput.Value()

	newPinnedList := lo.Map[list.Item, list.Item](b.singleList.Items(), func(i list.Item, _ int) list.Item {
		p := i.(list.PinnedListItem)
		p.SetCurrentValue(currentInputValue)
		return p
	})
	b.singleList.SetItems(newPinnedList)

	commands := b.generateCommandsFor(b.config.MultiLineConfigList, true)
	commands = append(commands, b.generateCommandsFor(b.config.SingleLineConfigList, false)...)

	commandsAsTeaCmds := lo.Map[shell.Command, tea.Cmd](commands, func(c shell.Command, _ int) tea.Cmd {
		return c.Run
	})

	teaCmds = append(teaCmds, commandsAsTeaCmds...)

	if len(currentInputValue) == 0 {
		b.multiList.UnfilterItems()
		b.multiList.Unselect()
	} else {
		b.multiList.SetFilterValue(currentInputValue)
		b.multiList.FilterItems()
		if len(b.multiList.VisibleItems()) > 0 {
			b.multiList.Select(0)
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
		if len(b.multiList.VisibleItems()) > 0 {
			autocompleteValue := b.multiList.VisibleItems()[0].FilterValue()
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

	sections = append(sections, b.singleList.View())
	sections = append(sections, queryStyle.Render(b.queryInput.View()))
	sections = append(sections, b.multiList.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (b *Bubble) focusSingleList() tea.Msg {
	b.multiList.Unselect()
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text
	b.singleList.Select(b.singleList.Height() - 1)

	return FocusChangeMsg{newFocus: SingleList}
}

func (b *Bubble) focusQueryInput() tea.Msg {
	b.queryInput.Focus()
	b.queryInput.PromptStyle = b.styles.QueryPromptFocused
	b.singleList.Unselect()
	b.multiList.Unselect()

	return FocusChangeMsg{newFocus: QueryInput}
}

func (b *Bubble) focusMultiList() tea.Msg {
	b.singleList.Unselect()
	b.multiList.Select(0)
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text

	return FocusChangeMsg{newFocus: MultiList}
}
