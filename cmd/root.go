/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"launcher/internal/config"
	"launcher/internal/tui"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pipe "github.com/b4b4r07/go-pipe"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

var cliCfg = config.LauncherConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "launcher",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		for _, sourceConfig := range cliCfg.SourceConfigList {
			pipedCommands := strings.Split(sourceConfig.Command,"|")
			execCommands := []*exec.Cmd{}
			for _, pipeCommand := range pipedCommands {
				commandComponents := strings.Split(strings.TrimSpace(pipeCommand)," ")
				mainCommand := commandComponents[0]
				args := commandComponents[1:]

				execCommands = append(execCommands, exec.Command(mainCommand, args...))
			}

			var b bytes.Buffer
			if err := pipe.Command(&b, execCommands...); err != nil {
				fmt.Println(err)
			}

			items := []list.Item{}
			lines := strings.Split(b.String(), "\n")
			for _, line := range lines {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					data := parts[0]
					label := parts[1]
					items = append(items, tui.NewLauncherListItem(label, data))
				}
			}

			list := list.New(items, list.NewDefaultDelegate(), 0, 0)
			var bubble *tui.Bubble = tui.NewBubble(&cliCfg, &list)
			p := tea.NewProgram(bubble, tea.WithAltScreen())

			if err := p.Start(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".launcher")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("There is no %HOME%/.launcher.yml present. Please create one.")
		} else {
			fmt.Println("Error reading config")
			fmt.Println(err)
		}
		os.Exit(1)
	}

	viper.Unmarshal(&cliCfg)
}
