package shell

import (
	"bytes"
	"os/exec"
	"strings"
	pipe "github.com/b4b4r07/go-pipe"
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

func (c Command) Run() (err error, buff bytes.Buffer) {
	var b bytes.Buffer
	return pipe.Command(&b, c.execCommands...), b
}

func (c Command) ResultLines() (e error, lines []string) {
	var b bytes.Buffer
	var err error

	err, b = c.Run()
	return err, strings.Split(b.String(), "\n")
}
