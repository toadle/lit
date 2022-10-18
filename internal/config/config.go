package config

import (
	"lit/internal/shell"
	"strings"

	"github.com/samber/lo"
)

type CommandGenerator struct {
	Command string `mapstructure:"command"`
}

type SourceConfig struct {
	CommandGenerator `mapstructure:",squash"`
	Action           string `mapstructure:"action"`
}

type SingleLineSourceConfig struct {
	SourceConfig `mapstructure:",squash"`
	Label        string `mapstructure:"label"`
}

type MultiLineSourceConfig struct {
	SourceConfig `mapstructure:",squash"`
	Format       string                `mapstructure:"format"`
	Labels       MultiLineLabelsConfig `mapstructure:"labels"`
}

type MultiLineLabelsConfig struct {
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
}

type LauncherConfig struct {
	SingleLineConfigList []SingleLineSourceConfig `mapstructure:"singleLine"`
	MultiLineConfigList  []MultiLineSourceConfig  `mapstructure:"multiLine"`
}

func (lc LauncherConfig) SingleLineSourceConfigFor(Command string) (SingleLineSourceConfig, bool) {
	return lo.Find(lc.SingleLineConfigList, func(sc SingleLineSourceConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) MultiLineSourceConfigFor(Command string) (MultiLineSourceConfig, bool) {
	return lo.Find(lc.MultiLineConfigList, func(sc MultiLineSourceConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) CommandGenerators() []CommandGenerator {
	commandGenerators := lo.Map(lc.MultiLineConfigList, func(mlc MultiLineSourceConfig, index int) CommandGenerator {
		return mlc.CommandGenerator
	})

	commandGenerators = append(commandGenerators, lo.Map(lc.SingleLineConfigList, func(slc SingleLineSourceConfig, index int) CommandGenerator {
		return slc.CommandGenerator
	})...)

	return commandGenerators
}

func (cg CommandGenerator) GenerateCommand(params map[string]string) *shell.Command {
	shellCmd := shell.NewCommand(cg.Command)
	if strings.Contains(cg.Command, "{input}") {
		shellCmd.SetParams(params)
	}
	return shellCmd
}
