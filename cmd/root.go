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
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/danish45007/gobackup/config"
	"github.com/spf13/cobra"
)

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
}
