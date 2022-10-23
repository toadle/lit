package cmd

import (
	"fmt"
	"os"

	"lit/internal/config"
	"lit/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliCfg = config.LauncherConfig{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lit",
	Short: "A command-line quicklauncher.",
	Long: `lit
An easily configurable quicklauncher for your terminal. 
You need to create a ~/.lit.yml config file in order to use it. 
See https://github.com/toadle/lit for help.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var bubble *tui.Bubble = tui.NewBubble(&cliCfg)
		p := tea.NewProgram(bubble)

		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
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
	viper.SetConfigName(".lit")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("There is no %HOME%/.lit.yml present. Please create one.")
		} else {
			fmt.Println("Error reading config")
			fmt.Println(err)
		}
		os.Exit(1)
	}

	viper.Unmarshal(&cliCfg)
}
