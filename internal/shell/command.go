package shell

import (
	"bytes"
	"os/exec"
	"strings"
	pipe "github.com/b4b4r07/go-pipe"
	tea "github.com/charmbracelet/bubbletea"
)

type Command struct {
	execCommands []*exec.Cmd
}

func NewCommand(str string) *Command {
	pipedCommands := strings.Split(str,"|")
	execCommands := []*exec.Cmd{}
	for _, pipeCommand := range pipedCommands {
		commandComponents := strings.Split(strings.TrimSpace(pipeCommand)," ")
		mainCommand := commandComponents[0]
		args := commandComponents[1:]

		execCommands = append(execCommands, exec.Command(mainCommand, args...))
	}
	return &Command{execCommands: execCommands}
}

func (c Command) Run() tea.Msg {
	var b bytes.Buffer

	if err := pipe.Command(&b, c.execCommands...); err != nil {
		return errMsg{err}
	}

	return ShellCommandResultMsg{Lines: strings.Split(b.String(), "\n")}
}

type ShellCommandResultMsg struct {
	Lines []string
}

type errMsg struct{ err error }
func (e errMsg) Error() string { return e.err.Error() }
