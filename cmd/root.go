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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config = LauncherConfig{}

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
		for _, sourceConfig := range config.SourceConfigList {
			
			commandComponents := strings.Split(sourceConfig.Command," ")
			mainCommand := commandComponents[0]
			args := commandComponents[1:]
			
			cmd := exec.Command(mainCommand, args...)
			cmd.Stdin = strings.NewReader("and old falcon")
	
			var out bytes.Buffer
			cmd.Stdout = &out
	
			err := cmd.Run()
	
			if err != nil {
				fmt.Println(err)
			}
			
			fmt.Print("%q\n", out.String())
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
	
	viper.Unmarshal(&config)
}