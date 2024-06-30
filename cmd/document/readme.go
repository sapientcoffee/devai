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

	"github.com/spf13/cobra"

	utils "buildey/pkg/common"
	genai "buildey/pkg/services"
)

var (
	codeFile   string
	promptFlag string
)

func generateReadme() {
	code, err := utils.ReadFileContents(codeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}
	chatPrompt := utils.BuildChatPrompt("Write documentation, in markdown, following DORA, stc.org, and writethedocs.org guidelines", promptFlag, code)

	fmt.Println("Generating and printing the release notes.")
	fmt.Println(genai.LangChainVertexChat(chatPrompt))

}

// readmeCmd represents the readme command
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "Create readme.md for the code passed in for documentation purposes.",
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
		promptFlag, err = cmd.Flags().GetString("prompt")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading prompt flag:", err)
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

		generateReadme()
	},
}

func init() {
	readmeCmd.Flags().StringP("git", "g", "", "The Git repo location to use for this action")
	readmeCmd.Flags().StringP("file", "f", "", "The file to use for this action")

	DocumentCmd.AddCommand(readmeCmd)
}
