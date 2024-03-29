package shell

import (
	"bytes"
	"os/exec"
	"strings"

	pipe "github.com/b4b4r07/go-pipe"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samber/lo"
)

type Command struct {
	cmdStr string
	params map[string]string
}

func NewCommand(str string) *Command {
	return &Command{cmdStr: str}
}

func (c Command) Run() tea.Msg {
	var b bytes.Buffer

	err := pipe.Command(&b, c.execCommands()...)

	return ShellCommandResultMsg{Output: b.String(), CmdStr: c.cmdStr, Successful: (err == nil)}
}

func (c *Command) SetParams(params map[string]string) {
	c.params = params
}

func (c Command) execCommands() []*exec.Cmd {
	pipedCommands := strings.Split(c.cmdStr, "|")

	return lo.Map(pipedCommands, func(pipeCommand string, _ int) *exec.Cmd {
		commandComponents := strings.Split(strings.TrimSpace(pipeCommand), " ")
		commandComponents = lo.Map(commandComponents, func(commandComponent string, _ int) string {
			return SetCommandParameters(commandComponent, c.params)
		})

		mainCommand := commandComponents[0]
		args := commandComponents[1:]

		return exec.Command(mainCommand, args...)
	})
}

type ShellCommandResultMsg struct {
	CmdStr     string
	Output     string
	Successful bool
}

func (s ShellCommandResultMsg) Lines() []string {
	return strings.Split(s.Output, "\n")
}
