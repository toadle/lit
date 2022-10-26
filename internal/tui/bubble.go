package tui

import (
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

type QueryChangedMsg int
type FocusChangeMsg struct {
	newFocus CursorFocus
}
type ResetMsg int
type CursorFocus int

const (
	CalculatorList CursorFocus = iota
	QueryInput
	SearchList
)

type Bubble struct {
	styles style.Styles
	config *config.LauncherConfig
	width  int
	height int

	searchList     list.Model
	queryInput     textinput.Model
	calculatorList list.Model

	queryInputTag int
	focus         CursorFocus
}

func NewBubble(cliCfg *config.LauncherConfig) *Bubble {
	b := &Bubble{
		styles:         style.DefaultStyles(),
		config:         cliCfg,
		searchList:     list.New([]list.Item{}, 0),
		queryInput:     textinput.New(),
		calculatorList: list.New([]list.Item{}, 0),
		focus:          QueryInput,
	}

	b.searchList.SetNoResultText("Nothing found.")
	b.queryInput.Placeholder = "Your Query"
	b.queryInput.ShowCompletions = true
	b.focusQueryInput()

	return b
}

func (b *Bubble) Init() tea.Cmd {
	var teaCmds []tea.Cmd

	for _, sourceConfig := range b.config.SearchConfigList {
		shellCmd := sourceConfig.GenerateCommand(map[string]string{
			"input": b.queryInput.Value(),
		})
		teaCmds = append(teaCmds, shellCmd.Run)
	}

	var calculatorItems []list.Item
	for _, sourceConfig := range b.config.CalculatorConfigList {
		calculatorItems = append(calculatorItems, list.NewCalculatorListItem(sourceConfig))
	}
	b.calculatorList.SetItems(calculatorItems)

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
			switch b.focus {
			case CalculatorList:
				i, ok := b.calculatorList.SelectedItem().(list.CalculatorListItem)
				if ok {
					params := map[string]string{
						"input": b.queryInput.Value(),
						"data":  i.Output(),
					}
					handler := b.generateEntrySelectedHandler(i.Action(), params)
					teaCmds = append(teaCmds, handler)
				}
			case QueryInput:
				if b.queryInput.CanBeCompleted() {
					availableCompletion := b.queryInput.AvailableCompletion()
					li, found := lo.Find[list.Item](b.searchList.VisibleItems(), func(listItem list.Item) bool {
						return listItem.FilterValue() == availableCompletion
					})

					if found {
						i, ok := li.(list.SearchListItem)
						if ok {
							handler := b.generateEntrySelectedHandler(i.Action(), i.Params())
							teaCmds = append(teaCmds, handler)
						}
					}
				}

			case SearchList:
				i, ok := b.searchList.SelectedItem().(list.SearchListItem)
				if ok {
					handler := b.generateEntrySelectedHandler(i.Action(), i.Params())
					teaCmds = append(teaCmds, handler)
				}
			}
		case tea.KeyUp:
			switch b.focus {
			case QueryInput:
				teaCmds = append(teaCmds, b.focusCalculatorList)
			case SearchList:
				if b.searchList.Index() == 0 {
					teaCmds = append(teaCmds, b.focusQueryInput)
				}
			}
		case tea.KeyDown:
			switch b.focus {
			case CalculatorList:
				if b.calculatorList.Index() == b.calculatorList.Height()-1 {
					teaCmds = append(teaCmds, b.focusQueryInput)
				}
			case QueryInput:
				teaCmds = append(teaCmds, b.focusSearchList)
			}
		case tea.KeyRunes, tea.KeyBackspace:
			if b.focus != QueryInput {
				teaCmds = append(teaCmds, b.focusQueryInput)

				//TODO: Find a better solution - This is a hack
				var cmd tea.Cmd
				b.queryInput.Focus()
				b.queryInput, cmd = b.queryInput.Update(msg)
				teaCmds = append(teaCmds, cmd)
				// hack ends here
			}
			b.queryInputTag++
			teaCmds = append(teaCmds, tea.Tick(time.Millisecond*100, func(_ time.Time) tea.Msg {
				return QueryChangedMsg(b.queryInputTag)
			}))
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		// _, right, _, left := styles.App.GetMargin()

		b.searchList.SetHeight(7)
		b.calculatorList.SetHeight(len(b.calculatorList.Items()))
	case shell.ShellCommandResultMsg:
		sourceConfig, isMultiLine := b.config.SearchConfigFor(msg.CmdStr)
		if isMultiLine {
			items := b.searchList.Items()
			for _, line := range msg.Lines() {
				items = append(items, list.NewSearchListItem(line, sourceConfig))
			}
			b.searchList.SetItems(items)
		} else {
			newPinnedList := lo.Map[list.Item, list.Item](b.calculatorList.Items(), func(i list.Item, _ int) list.Item {
				p := i.(list.CalculatorListItem)
				if p.CmdStr() == msg.CmdStr {
					p.SetOutput(msg.Output)
					p.SetSuccessful(msg.Successful)
				}
				return p
			})
			b.calculatorList.SetItems(newPinnedList)
		}

	case QueryChangedMsg:
		if int(msg) == b.queryInputTag {
			teaCmds = append(teaCmds, b.handleQueryChanged()...)
		}
	case FocusChangeMsg:
		b.focus = FocusChangeMsg(msg).newFocus

	case ResetMsg:
		b.queryInput.SetValue("")
		teaCmds = append(teaCmds, b.focusQueryInput)
	}

	var cmd tea.Cmd

	if b.focus == QueryInput {
		b.queryInput, cmd = b.queryInput.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == SearchList {
		b.searchList, cmd = b.searchList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	if b.focus == CalculatorList {
		b.calculatorList, cmd = b.calculatorList.Update(msg)
		teaCmds = append(teaCmds, cmd)
	}

	return b, tea.Batch(teaCmds...)
}

func (b *Bubble) generateEntrySelectedHandler(action string, params map[string]string) tea.Cmd {
	return func() tea.Msg {
		shellCmd := shell.NewCommand(action)
		shellCmd.SetParams(params)
		shellCmd.Run()

		if b.config.CloseOnAction {
			return tea.Quit()
		} else {
			return ResetMsg(0)
		}

	}
}

func (b *Bubble) handleQueryChanged() []tea.Cmd {
	var teaCmds []tea.Cmd

	currentInputValue := b.queryInput.Value()

	newPinnedList := lo.Map(b.calculatorList.Items(), func(i list.Item, _ int) list.Item {
		p := i.(list.CalculatorListItem)
		p.SetCurrentValue(currentInputValue)
		return p
	})
	b.calculatorList.SetItems(newPinnedList)

	commands := b.generateCommandsFor(b.config.CommandGenerators())

	commandsAsTeaCmds := lo.Map(commands, func(c shell.Command, _ int) tea.Cmd {
		return c.Run
	})

	teaCmds = append(teaCmds, commandsAsTeaCmds...)

	if len(currentInputValue) == 0 {
		b.searchList.UnfilterItems()
		b.searchList.Unselect()
	} else {
		b.searchList.SetFilterValue(currentInputValue)
		b.searchList.FilterItems()
	}

	teaCmds = append(teaCmds, b.generateCompletions)
	return teaCmds
}

func (b Bubble) generateCommandsFor(configs []config.CommandGenerator) []shell.Command {
	currentInputValue := b.queryInput.Value()
	shellCommands := []shell.Command{}
	for _, cg := range configs {
		if strings.Contains(cg.Command, "{input}") {
			shellCmd := shell.NewCommand(cg.Command)
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
		if len(b.searchList.VisibleItems()) > 0 {
			autocompleteValue := b.searchList.VisibleItems()[0].FilterValue()
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

	sections = append(sections, b.calculatorList.View())
	sections = append(sections, queryStyle.Render(b.queryInput.View()))
	sections = append(sections, b.searchList.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (b *Bubble) focusCalculatorList() tea.Msg {
	b.searchList.Unselect()
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text
	b.calculatorList.Select(b.calculatorList.Height() - 1)

	return FocusChangeMsg{newFocus: CalculatorList}
}

func (b *Bubble) focusQueryInput() tea.Msg {
	b.queryInput.Focus()
	b.queryInput.PromptStyle = b.styles.QueryPromptFocused
	b.calculatorList.Unselect()
	b.searchList.Unselect()

	return FocusChangeMsg{newFocus: QueryInput}
}

func (b *Bubble) focusSearchList() tea.Msg {
	b.calculatorList.Unselect()
	b.searchList.Select(0)
	b.queryInput.Blur()
	b.queryInput.PromptStyle = b.styles.Text

	return FocusChangeMsg{newFocus: SearchList}
}
