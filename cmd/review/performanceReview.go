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
package review

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	utils "buildey/pkg/common"
	genaiService "buildey/pkg/services"
)

func performanceReview() {
	code, err := utils.ReadFileContents(codeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	chatPrompt := utils.BuildChatPrompt(
		`Please optimize and specifically describe the changes you would make
		`,
		promptFlag, code)

	fmt.Println("Generating and printing the performance review. Output in valid json format.")

	// Select which module will be used to call vertex, options are LangChainVertexChat or VertexChatText which uses vertexai SDK
	// response := genaiService.LangChainVertexChat(chatPrompt)
	response, err := genaiService.VertexChatText(chatPrompt)
	if err != nil {
		fmt.Printf("Error calling VertexChatText: %v\n", err)
		return
	}

	// Parse JSON response
	// var reviewItems []utils.ReviewItem
	// err = json.Unmarshal([]byte(response), &reviewItems) // Unmarshall directly into reviewItems
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
	// 	return
	// }

	// tableContents := utils.CreateTable(reviewItems, "markdown")
	// tableContents := utils.CreateTable(reviewItems, "cli")

	// fmt.Println(tableContents)
	fmt.Println(response)

}

// performanceReviewCmd represents the performanceReview command
var performanceReviewCmd = &cobra.Command{
	Use:   "performance",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("performanceReview called")

		fileFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file flag:", err)
			return
		}

		if fileFlag != "" {
			codeFile = fileFlag
		} else {
			fmt.Println("No file specified, using default.")
			codeFile = "../example_code/coffee.go" // Consider moving default to a constant
		}
		fmt.Println("Performance Review for " + fileFlag + " .....")
		performanceReview()
	},
}

func init() {

	performanceReviewCmd.Flags().StringP("file", "f", "", "The file containing the code to be reviewed")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// performanceReviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// performanceReviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
