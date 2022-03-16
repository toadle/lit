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

type item struct {
		title, desc string
	}

	func (i item) Title() string       { return i.title }
	func (i item) Description() string { return i.desc }
	func (i item) FilterValue() string { return i.title }


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
			//fmt.Print(b.String())
			items := []list.Item{
				item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
				item{title: "Nutella", desc: "It's good on toast"},
				item{title: "Bitter melon", desc: "It cools you down"},
				item{title: "Nice socks", desc: "And by that I mean socks without holes"},
				item{title: "Eight hours of sleep", desc: "I had this once"},
				item{title: "Cats", desc: "Usually"},
				item{title: "Plantasia, the album", desc: "My plants love it too"},
				item{title: "Pour over coffee", desc: "It takes forever to make though"},
				item{title: "VR", desc: "Virtual reality...what is there to say?"},
				item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
				item{title: "Linux", desc: "Pretty much the best OS"},
				item{title: "Business school", desc: "Just kidding"},
				item{title: "Pottery", desc: "Wet clay is a great feeling"},
				item{title: "Shampoo", desc: "Nothing like clean hair"},
				item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
				item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
				item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
				item{title: "Stickers", desc: "The thicker the vinyl the better"},
				item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
				item{title: "Warm light", desc: "Like around 2700 Kelvin"},
				item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
				item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
				item{title: "Terrycloth", desc: "In other words, towel fabric"},
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
