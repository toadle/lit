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

type CalculatorConfig struct {
	SourceConfig `mapstructure:",squash"`
	Label        string `mapstructure:"label"`
}

type SearchConfig struct {
	SourceConfig `mapstructure:",squash"`
	Format       string                `mapstructure:"format"`
	Labels       MultiLineLabelsConfig `mapstructure:"labels"`
}

type MultiLineLabelsConfig struct {
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
}

type LauncherConfig struct {
	CalculatorConfigList []CalculatorConfig `mapstructure:"calculators"`
	SearchConfigList     []SearchConfig     `mapstructure:"searches"`
}

func (lc LauncherConfig) CalculatorConfigFor(Command string) (CalculatorConfig, bool) {
	return lo.Find(lc.CalculatorConfigList, func(sc CalculatorConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) SearchConfigFor(Command string) (SearchConfig, bool) {
	return lo.Find(lc.SearchConfigList, func(sc SearchConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) CommandGenerators() []CommandGenerator {
	commandGenerators := lo.Map(lc.SearchConfigList, func(mlc SearchConfig, index int) CommandGenerator {
		return mlc.CommandGenerator
	})

	commandGenerators = append(commandGenerators, lo.Map(lc.CalculatorConfigList, func(slc CalculatorConfig, index int) CommandGenerator {
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
