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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	utils "buildey/pkg/common"
	genaiService "buildey/pkg/services"
)

var (
	codeFile   string
	promptFlag string
)

// codeReview performs the code review process.
func codeReview() {
	code, err := utils.ReadFileContents(codeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return
	}

	chatPrompt := utils.BuildChatPrompt(
		`### instruction ###

		You are an enterprise developer for cymbal coffee. You are an expert in software development with over 20 years experience. There is a codebase in "current data context" section which you should do the review on. Each question should be review independently against the "current data context".
		
		Ensure that the question is relevant to the code snippet (codebase shared), if it does not look applicable reply with "n/a".
		
		First questions is to detect violations of coding style guidelines and conventions. Identify inconsistent formatting, naming conventions, indentation, comment placement, and other style-related issues. Provide suggestions to the detected violations to maintain a consistent and readable codebase if this is a problem.
		
		Second questions is to identify common issues such as code smells, anti-patterns, potential bugs, performance bottlenecks, and security vulnerabilities. Offer actionable recommendations to address these issues and improve the overall quality of the code.
		
		
		### example diagogs ###
		<query> First questions are to detect violations of coding style guidelines and conventions. Identify inconsistent formatting, naming conventions, indentation, comment placement, and other style-related issues. Provide suggestions or automatically fix the detected violations to maintain a consistent and readable codebase if this is a problem.
		import "fmt"
		
		func main() {
			name := "Alice"
			greeting := fmt.Sprintf("Hello, %s!", name)
			fmt.Println(greeting)
		}
		
		
		<response> [
			{
				"question": "Indentation",
				"answer": "yes",
				"description": "Code is consistently indented with spaces (as recommended by Effective Go)"
			},
			{
				"question": "Variable Naming",
				"answer": "yes",
				"description": "Variables ("name", "greeting") use camelCase as recommended"
			},
			{
				"question": "Line Length",
				"answer": "yes",
				"description": "Lines are within reasonable limits" 
			},
			{
				"question": "Package Comments", 
				"answer": "n/a",
				"description": "This code snippet is too small for a package-level comment"
			}
		]
		
		
		<query> Identify common issues such as code smells, anti-patterns, potential bugs, performance bottlenecks, and security vulnerabilities. Offer actionable recommendations to address these issues and improve the overall quality of the code.
		
		"package main
		
		import (
			"fmt"
			"math/rand"
			"time"
		)
		
		// Global variable, potentially unnecessary 
		var globalCounter int = 0 
		
		func main() {
			items := []string{"apple", "banana", "orange"}
		
			// Very inefficient loop with nested loop for a simple search
			for _, item := range items {
				for _, search := range items {
					if item == search {
						fmt.Println("Found:", item)
					}
				}
			}
		
			// Sleep without clear reason, potential performance bottleneck
			time.Sleep(5 * time.Second) 
		
			calculateAndPrint(10)
		}
		
		// Potential divide-by-zero risk
		func calculateAndPrint(input int) {
			result := 100 / input 
			fmt.Println(result)
		}"
		
		<response> [
			{
				"question": "Global Variables",
				"answer": "no",
				"description": "Potential issue: Unnecessary use of the global variable 'globalCounter'. Consider passing values as arguments for better encapsulation." 
			},
			{
				"question": "Algorithm Efficiency",
				"answer": "no",
				"description": "Highly inefficient search algorithm with an O(n^2) complexity. Consider using a map or a linear search for better performance, especially for larger datasets."
			},
			{
				"question": "Performance Bottlenecks",
				"answer": "no",
				"description": "'time.Sleep' without justification introduces a potential performance slowdown. Remove it if the delay is unnecessary or provide context for its use."
			},
			{
				"question": "Potential Bugs",
				"answer": "no",
				"description": "'calculateAndPrint' function has a divide-by-zero risk. Implement a check to prevent division by zero and handle the error appropriately."
			},
			{ 
				"question": "Code Readability",
				"answer": "no",
				"description": "Lack of comments hinders maintainability. Add comments to explain the purpose of functions and blocks of code."
			} 
		]
		
		### output details ####
		
		Create a JSON output which provides a response to each of the questions. The output should be in the format of a JSON array with each element containing - question, answer, description. The answer should only be "yes" if it is acceptable, "no" if it has problems or recommendations or if the question is not relevant return "n/a". Any other detail should be provided in the description field.
		
		### current data context ###
		`,
		promptFlag, code)

	// fmt.Println(chatPrompt)
	fmt.Println("Generating and printing the code review.")

	response, err := genaiService.GetAIResponse(chatPrompt, "")
	// response, err := genaiService.GetAIResponse(chatPrompt, "VertexEngine") //VertexEngine or LangChainEngine or CodeyEngine

	// response, err := genaiService.GetAIResponse(chatPrompt)
	if err != nil {
		fmt.Printf("Error getting AI response: %v\n", err)
		return
	}

	validatedJSON, err := utils.ValidateJSON(response)
	if err != nil {
		fmt.Printf("error validating JSON:  %v\n", err)
		return
	}

	// Parse JSON response
	var reviewItems []utils.ReviewItem
	err = json.Unmarshal([]byte(validatedJSON), &reviewItems) // Unmarshall directly into reviewItems
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON response: %v\n", err)
		return
	}

	// tableContents := utils.CreateTable(reviewItems, "markdown")
	tableContents := utils.CreateTable(reviewItems, "cli")

	fmt.Println(tableContents)
	// fmt.Println(response) // Print the JSON response

}

// codeReviewCmd represents the codeReview command
var codeReviewCmd = &cobra.Command{
	Use:   "code",
	Short: "Assist with code reviews",
	Long:  `Assists with code reviews and generates a table with AI review results. The goal is to help with the "heavy" lifting of code reviews to allow the reviewer to focus on adding value.`,
	Run: func(cmd *cobra.Command, args []string) {

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
		fmt.Println("Code Review for " + fileFlag + " .....")
		codeReview()
	},
}

func init() {
	codeReviewCmd.Flags().StringP("file", "f", "", "The file containing the code to be reviewed")
}
