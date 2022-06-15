package config

import (
	"github.com/samber/lo"
)

type SourceConfig struct {
	Command string `mapstructure:"command"`
	ItemFormat string `mapstructure:"itemFormat"`
	WhenSelected string `mapstructure:"whenSelected"`
	Pinned bool `mapstructure:"pinned"`
	Label string `mapstructure:"label"`
}

type LauncherConfig struct {
	SourceConfigList []SourceConfig `mapstructure:"sources"`
}

func (lc LauncherConfig) SourceConfigFor(Command string) (SourceConfig, bool) {
	return lo.Find[SourceConfig](lc.SourceConfigList, func(sc SourceConfig) bool {
		return Command == sc.Command
	})
}

func (lc LauncherConfig) ResultSourceConfigList() []SourceConfig {
	return lo.Filter[SourceConfig](lc.SourceConfigList, func(sc SourceConfig, _ int) bool {
		return !sc.Pinned
	})
}

func (lc LauncherConfig) PinnedSourceConfigList() []SourceConfig {
	return lo.Filter[SourceConfig](lc.SourceConfigList, func(sc SourceConfig, _ int) bool {
		return sc.Pinned
	})
}
