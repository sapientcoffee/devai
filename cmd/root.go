/*
Copyright 2024 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

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

	"buildey/cmd/document"
	"buildey/cmd/info"
	"buildey/cmd/release"
	"buildey/cmd/review"
	"buildey/cmd/inspect"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "20240414"
	Verbose bool
	Debug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "buildey",
	Version: version,
	Short:   "A GenAI assistent for usage developers in CLI form",
	Long: `A CLI to leverage GenAI to help with common developer activities.

buildey could be used standalone or embedded into CI workflows to assist 
humans and do heavy lifting to improve productivity and help with all 
those things we know are important.
Thinks documentation, code reviews, QA, performance etc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}
		}
		if Verbose {
			log.Info("Verbose flag NOT Implemented yet")
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(review.ReviewCmd)
	rootCmd.AddCommand(release.ReleaseCmd)
	rootCmd.AddCommand(document.DocumentCmd)
	rootCmd.AddCommand(inspect.InspectCmd)
	rootCmd.AddCommand(info.InfoCmd)

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.buildey.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

}
