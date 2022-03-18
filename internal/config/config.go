package config

type SourceConfig struct {
    Command string `mapstructure:"command"`
    ItemFormat string `mapstructure:"itemFormat"`
    WhenSelected string `mapstructure:"whenSelected"`
}

type LauncherConfig struct {
    SourceConfigList []SourceConfig `mapstructure:"sources"`
}
