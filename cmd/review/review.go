/*
Copyright © 2024 Google LLC

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
package review

import (
	"github.com/spf13/cobra"
)

// var (
// 	file string
// )

// reviewCmd represents the review command
var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "A number of options to assist with reviews of certain aspects of code",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ReviewCmd.AddCommand(codeReviewCmd)
	ReviewCmd.AddCommand(performanceReviewCmd)
	ReviewCmd.AddCommand(archReviewCmd)

	codeReviewCmd.Flags().StringP("git", "g", "", "The Git repo location to use for this action")

	// codeReviewCmd.Flags().StringVarP(&file, "file", "f", "./diff.txt", "Filename with the code.")

}
