package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"launcher/internal/tui/style"
	"launcher/internal/config"
)

type Bubble struct {
	config	*config.LauncherConfig
	list	*list.Model
}

func NewBubble(cliCfg *config.LauncherConfig, list *list.Model) *Bubble {
    b := &Bubble{
      config:      cliCfg,
	  list:		   list,
    }
    return b
  }

func (b *Bubble) Init() tea.Cmd {
  return nil
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    if msg.String() == "ctrl+c" {
      return b, tea.Quit
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
