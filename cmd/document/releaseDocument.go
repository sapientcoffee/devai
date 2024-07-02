/*
Copyright 2023 Google LLC

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

package document

import (
	"fmt"
	"os"

	utils "buildey/pkg/common"
	genaiService "buildey/pkg/services"

	"github.com/spf13/cobra"
)

func generateReleaseNotes() {
	code, err := utils.ReadFileContents(codeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	chatPrompt := utils.BuildChatPrompt("Create detailed release notes", promptFlag, code)

	fmt.Println("Generating and printing the release notes.")

	fmt.Println(genaiService.GetAIResponse(chatPrompt, ""))

}

var releaseCmd = &cobra.Command{
	Use:   "release notes",
	Short: "Create release notes for documentation purposes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitFlag, err := cmd.Flags().GetString("git")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading git flag:", err)
			return
		}
		fileFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file flag:", err)
			return
		}

		if gitFlag != "" {
			fmt.Println("Git flag is not yet implemented:", gitFlag)
			return
		}

		if fileFlag != "" {
			codeFile = fileFlag
		} else {
			fmt.Println("No file specified, using default.")
			codeFile = "../example_code/coffee.go" // Consider moving default to a constant
		}
		fmt.Println("Release Notes for " + fileFlag + " .....")
		generateReleaseNotes()
	},
}

func init() {
	releaseCmd.Flags().StringP("git", "g", "", "The Git repo location to use for this action")
	releaseCmd.Flags().StringP("file", "f", "", "The file to use for this action")

	DocumentCmd.AddCommand(releaseCmd)
}
