package tui

import (
	"strings"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"launcher/internal/tui/style"
	"launcher/internal/config"
)

type LauncherListItem struct {
	label, data, itemFormat, whenSelected string
}

func (i LauncherListItem) Title() string {
	return i.label
}
func (i LauncherListItem) Description() string {
	return i.data
}
func (i LauncherListItem) FilterValue() string {
	return i.label
}

func NewLauncherListItem(itemData, itemFormat, whenSelected string) *LauncherListItem {
	var label, data string
	parts := strings.Split(itemData, ":")
	if len(parts) == 2 {
		data = parts[0]
		label = parts[1]
	}

	return &LauncherListItem{label: label, data: data, itemFormat: itemFormat, whenSelected: whenSelected}
}

type Bubble struct {
	config	*config.LauncherConfig
	list	*list.Model
}

func NewBubble(cliCfg *config.LauncherConfig, list *list.Model) *Bubble {
	b := &Bubble{
		config:	cliCfg,
		list:	list,
	}
	return b
}

func (b *Bubble) Init() tea.Cmd {
	return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return b, tea.Quit
		case "enter", " ":
			i, ok := b.list.SelectedItem().(*LauncherListItem)
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
		top, right, bottom, left := style.DefaultStyles().App.GetMargin()
		b.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	var list list.Model
	list, cmd = b.list.Update(msg)
	b.list = &list
	return b, cmd
}

func (b *Bubble) View() string {
	return style.DefaultStyles().App.Render(b.list.View())
}
