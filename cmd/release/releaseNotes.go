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
package release

import (
	genai "buildey/pkg/services"
	"fmt"

	"github.com/spf13/cobra"
)

func generateReadme(codeSnippet string) (string, error) {
	llmPrompt := fmt.Sprintf("Explain the following Go code snippet in markdown, following DORA, stc.org, and writethedocs.org guidelines:\n\n%s", codeSnippet)
	releaseNotes := genai.LangChainVertexChat(llmPrompt)
	return releaseNotes, nil // Return releaseNotes and nil error
}

// releaseNotesCmd represents the releaseNotes command
var releaseNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Generate release notes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Release Note Generate for ")

		codeSnippet := "package main\n// ... (rest of the code)" // Or read from a file
		releaseNotes, err := generateReadme(codeSnippet)
		if err != nil {
			fmt.Println("Error generating release notes:", err)
			return
		}
		fmt.Println(releaseNotes)

	},
}

func init() {
	ReleaseCmd.AddCommand(releaseNotesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseNotesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseNotesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
