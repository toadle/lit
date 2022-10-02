package config

import (
	"github.com/samber/lo"
)

type SourceConfig struct {
	Command string `mapstructure:"command"`
	Format  string `mapstructure:"format"`
	Action  string `mapstructure:"action"`
	Label   string `mapstructure:"label"`
}

type LauncherConfig struct {
	SingleLineConfigList []SourceConfig `mapstructure:"singleLine"`
	MultiLineConfigList  []SourceConfig `mapstructure:"multiLine"`
}

func (lc LauncherConfig) SingleLineSourceConfigFor(Command string) (SourceConfig, bool) {
	return lo.Find[SourceConfig](lc.SingleLineConfigList, func(sc SourceConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) MultiLineSourceConfigFor(Command string) (SourceConfig, bool) {
	return lo.Find[SourceConfig](lc.MultiLineConfigList, func(sc SourceConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) SourceConfigFor(Command string, multiline bool) (SourceConfig, bool) {
	if multiline {
		return lc.MultiLineSourceConfigFor(Command)
	} else {
		return lc.SingleLineSourceConfigFor(Command)
	}
}
