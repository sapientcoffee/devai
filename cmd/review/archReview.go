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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	utils "buildey/pkg/common"
	genaiService "buildey/pkg/services"
)

var (
	code string

// promptFlag string
)

// archReview performs the architecture review process.
func archReview() {
	// Create a string builder to accumulate code from files
	var codeBuilder strings.Builder

	// Walk the directory recursively
	err := filepath.Walk(code, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories themselves
		if !info.IsDir() {
			// Read file contents
			fileCode, err := utils.ReadFileContents(path)
			if err != nil {
				return err
			}

			// Append file contents and path to the builder
			codeBuilder.WriteString(fmt.Sprintf("## Code from file: %s\n", path))
			codeBuilder.WriteString(fileCode)
			codeBuilder.WriteString("\n\n")
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading directory:", err)
		return
	}

	// Get the accumulated code from the string builder
	accumulatedCode := codeBuilder.String()

	// --- Display the aggregated code ---
	fmt.Println("----- Aggregated Code -----")
	fmt.Println(accumulatedCode)
	fmt.Println("----- End of Aggregated Code -----")

	chatPrompt := utils.BuildChatPrompt(
		`### instruction ###

		You are an expert in Google Cloud Architectures and only use infromation provided in the pages contained under https://cloud.google.com/architecture/framework.
        Your goal is to compare the deployment code to the  Google Cloud Architecture Framework and recommend improvements to the code and design documents. 
        Structure your answer in the five sections of the framework: Operational Excellence, Security, Reliability, Cost Optimisation, and Performance Optimisation. 
        Your answer should focus on specific code or infrastructure improvements.
		
		### output details ####
		
		The output should be in markdown format with the highlevel sections that are reflected in the Google Cloud Architecture Framework pages as 2nd level titles. Also provide a date and timestamp of the when the latest update to the framework was made - in the overview page of each section it is found in the format "Last reviewed 2024-03-29 UTC".

		Evaluate the code, and recommend improvements in the code base, based on the recommendations of the Google Cloud Architecture Framework.If there points where the code does already adhere to the recommendations of the framework, breifly highlight those",		


		### example diagogs ###
		
		# Cloud Architecure Review
		Here are the recommendations of the Google Cloud Architecture framework:
		** Overview: **
		Reviewed against GCAF updated on: <"Last reviewed">

		In general the provided code compared against the architecture framework has the following possitives and areas that we reccomend invertigating for potential imrpovement.


		Operational Excellence:
		Reviewed against GCAF updated on: <"Last reviewed">

		Security, Privacy and compliance:
		Reviewed against GCAF updated on: <"Last reviewed">

		Reliability:
		Reviewed against GCAF updated on: <"Last reviewed">

		Cost Optimisation: 
		Reviewed against GCAF updated on: <"Last reviewed">
		
		Performance Optimisation:  
		Reviewed against GCAF updated on: <"Last reviewed">

		Sustainability:
		Reviewed against GCAF updated on: <"Last reviewed">

		### current data context ###
		`,
		promptFlag, accumulatedCode) // Use accumulatedCode here

	// fmt.Println(chatPrompt)
	fmt.Println("Generating and printing the code review.")

	response, err := genaiService.GetAIResponse(chatPrompt, "")
	// response, err := genaiService.GetAIResponse(chatPrompt, "VertexEngine") //VertexEngine or LangChainEngine or CodeyEngine

	// response, err := genaiService.GetAIResponse(chatPrompt)
	if err != nil {
		fmt.Printf("Error getting AI response: %v\n", err)
		return
	}

	fmt.Println(response)
	// fmt.Println(response) // Print the JSON response

}

// archReviewCmd represents the archReview command
var archReviewCmd = &cobra.Command{
	Use:   "arch",
	Short: "Assist with architecture reviews",
	Long:  `Assists with architecture reviews and generates a table with AI review results. The goal is to help with the "heavy" lifting of code reviews to allow the reviewer to focus on adding value.`,
	Run: func(cmd *cobra.Command, args []string) {

		dirFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file flag:", err)
			return
		}

		if dirFlag != "" {
			code = dirFlag
		} else {
			fmt.Println("No directory specified, using current directory.")
			code = "." // Default to the current directory
		}
		fmt.Println("Architecture Review for " + dirFlag + " .....")
		archReview()
	},
}

func init() {
	archReviewCmd.Flags().StringP("file", "f", "", "The directory containing the code to be reviewed")
}
