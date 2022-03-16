package config

type SourceConfig struct {
    Command string `mapstructure:"command"`
}

type LauncherConfig struct {
    SourceConfigList []SourceConfig `mapstructure:"sources"`
}