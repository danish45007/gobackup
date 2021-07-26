/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/danish45007/gobackup/config"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/itrepablik/itrlog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppInfo struct {
	Name    string
	Version string
}

type ConfigCommandDir struct {
	Ignore []string
	Log    bool
}

var cfgFile string

var appConfig = []AppInfo{}
var conCommand = []ConfigCommandDir{}
var ingoreDir []string
var logDir bool = false

var logTimeFormat = "Jan 02 2001 00:00:00 PM"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Short:   config.AppShortDesc,
	Long:    config.AppLongtDesc,
	Version: config.AppVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	myFigure := figure.NewColorFigure(config.DisplayName, "rectangles", "green", true)
	myFigure.Print()
	fmt.Println()
	LoadViperConfig()
	viper.WatchConfig() // Tell the viper to watch any new changes to the config file.
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		itrlog.Warn("Config file changed:", e.Name)
		LoadViperConfig()
	})
}

func LoadViperConfig() {
	cobra.OnInitialize(initConfig)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	// Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create the "config.yaml" asap.
			f, err := os.OpenFile("config.yaml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				itrlog.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
		} else {
			// Config file was found but another error was produced
			itrlog.Fatalf("fatal error config file: %v", err)
		}
	}

	err := viper.UnmarshalKey("app", &appConfig)
	if err != nil {
		itrlog.Error(err)
	}

	for _, a := range appConfig {
		if len(strings.TrimSpace(a.Name)) == 0 || strings.TrimSpace(a.Name) != config.AppName {
			errMsg := "app name must be " + config.AppName + ", please follow this naming convention to your config.yaml file under 'app'."
			color.Red(errMsg)
			itrlog.Fatal(errors.New(errMsg))
		}
		if len(strings.TrimSpace(a.Version)) == 0 || strings.TrimSpace(a.Version) != config.AppVersion {
			errMsg := "app version must be " + config.AppVersion + ", please follow this naming convention to your config.yaml file under 'app'."
			color.Red(errMsg)
			itrlog.Fatal(errors.New(errMsg))
		}
	}

	error := viper.UnmarshalKey("default.command_properties.comdir", &conCommand)
	if error != nil {
		itrlog.Error(error)
	}

	for _, data := range conCommand {
		ingoreDir = data.Ignore
		logDir = data.Log
	}

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.

		home, err := homedir.Dir()
		if err != nil {
			color.Red(err.Error())
			itrlog.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gobackup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gobackup")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		color.Red("Using config file: " + viper.ConfigFileUsed())
	}
}
