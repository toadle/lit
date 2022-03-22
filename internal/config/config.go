package config

import "github.com/thoas/go-funk"

type SourceConfig struct {
	Command string `mapstructure:"command"`
	ItemFormat string `mapstructure:"itemFormat"`
	WhenSelected string `mapstructure:"whenSelected"`
}

type LauncherConfig struct {
	SourceConfigList []SourceConfig `mapstructure:"sources"`
}

func (lc LauncherConfig) SourceConfigFor(Command string) SourceConfig {
	return funk.Find(lc.SourceConfigList, func(sc SourceConfig) bool {
		return Command == sc.Command
	}).(SourceConfig)
}
