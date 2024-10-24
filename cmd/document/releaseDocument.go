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

	chatPrompt := utils.BuildChatPrompt(
		`### instruction ###
		You are a technical writer creating release notes. 

		The audience for the release notes consists of experienced developers who integrate or add features to the application and service. They rely on detailed documentation, advance notice of upcoming features, and clear information about potential breaking changes.
		
		SOURCE CONTENT IS A FILE DIFF:

		You will be provided with the diff output of all changed files in the release. The diff describes the changes in the reference documentation, which directly reflects how the code has changed. Your task is to analyze the diff and clearly describe the changes. I will pull from your descriptions to populate the release notes.

		IGNORE GRAMMAR/STYLE CHANGES

		The diff is comprehensive and includes many grammar and style changes to existing definitions. Ignore these, as they are minor cosmetic updates for readability that don’t need to be included in the release notes.

		INTERPRETING THE DIFF SYNTAX:

		The diff output uses the following syntax to indicate changes:

		+ : This symbol indicates a line that was added in the new version.
		- : This symbol indicates a line that was removed in the new version.
		@@ ... @@ : These lines show the line numbers where changes occur in each file.
		The first line number refers to the original file (before changes).
		The second line number refers to the new file (after changes).
		<del>: This tag is used to indicate deprecated code. It is often used around method or class names in the documentation.
		
		You should focus on changes that affect functionality or integrations. Ignore internal implementation details, minor comment updates, or stylistic changes.
		Pay close attention to the <del> tag to identify deprecated elements.
		
		Analyze the diff and describe the changes that have been made using plain, readable language. Your analysis will mostly be matter of fact, describing the changes. The file diff will not tell you why the changes have been made or what the larger purpose is behind the changes – that is all right, as I will supplement the matter-of-fact changes with this larger context from other sources. Your task is mainly to describe the differences in the file diff. Especially not the following:

			New features: Describe any added classes, methods, or capabilities. Extrapolate the descriptions and purposes for the elements from the code.
			Deprecations: Identify any deprecated classes, methods, or fields.
			Other changes: Report significant changes that could affect functionality or integrations. Ignore internal implementation details, minor comment updates, or stylistic changes that don’t affect functionality.

		You can note any other significant changes too.

		### example ###
		<example 1>
		--- a/file.java
		+++ b/file.java
		@@ -1,5 +1,6 @@
		
		public class MyClass {
		-  private int value = 10;
		+  private int value = 20;  // This line was changed
		+  public void newValue() { ... } 
		}
		Explanation:
			The line private int value = 10; was removed.
			The line private int value = 20; was added (and is also marked with a comment).
			A new method public void newValue() { ... } was added.

		<example 2>
		EXAMPLE RELEASE NOTE SNIPPETS:

			New feature:

			### Improved data processing speed

			The "processData()" method now includes a "report_fields" parameter, which lets you specify a report for the list of data you want to process.
			Deprecation:

			### Data processing deprecations

			**"DataProcessor.Builder" class:**

			* **"setOldAlgorithm()"**: This method, used to configure the old processing algorithm, has been deprecated. Use the "setNewAlgorithm()" method instead. The old algorithm will be removed in a future release. 
			Removal:

			### Report processing deprecations

			**"ReportProcessor.Builder" class:**

			* **"setOldReport()"**: This deprecated method, used to configure the old processing algorithm, has been removed from the API. Use the "setNewReport()" method instead.
			Documentation update:


			### Documentation updates

			* The documentation for the "DataProcessor" class has been updated to include a new section on performance optimization. 

		### output details ####
		STYLE RULES:
			* Use sentence-case capitalization for headings and descriptions. In other words, only capitalize the first word in headings and subheadings.
			* Avoid adjectives. Be plain and clear. This is technical documentation, not marketing material.
			* Be precise with field and class names, using correct capitalization and the full name (for example, DataProcessor.Builder, calculateValue()).
			* Keep the language concise and technical, targeting experienced developers.

		MARKDOWN SYNTAX:
			* Provide the output in markdown formatting.
			* For the title, use the header # Release.
			* Use markdown headers (##, ###, ####) to structure the subheadings.
			* Surround classes, methods, and fields in backticks.
			* Speak in second-person voice (“you”) directly to developers.
		
		### current data context ###
		Here is the diff output for changes in this release:
		`,
		promptFlag, code)

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
