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

// Package genaiService provides functions for interacting with the Google Cloud VertexAI Codey API.

package genaiService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

// VertexChatText uses the Google Cloud VertexAI Codey API to generate text responses to user prompts.
//
// Args:
//
//	chatPrompt: The user's prompt to be sent to the VertexAI model.
//
// Returns:
//
//	A string containing the generated text response from the VertexAI model.
//	An error if any occurred during the process.
func VertexChatText(chatPrompt string) (string, error) {

	// Project ID and region for the VertexAI service.

	// project := os.Getenv("PROJECT_ID")
	// location := os.Getenv("REGION")
	project := "coffee-and-codey"
	region := "us-central1"

	// Choose the desired VertexAI model.
	// Options:
	//   - gemini-1.5-flash-001: A fast and efficient model.
	//   - gemini-1.5-pro: A more powerful model with better performance.
	//   - gemini-pro: The most powerful model with the best performance.
	model := "gemini-1.5-pro-002"

	// model := "gemini-1.5-flash-001"
	fmt.Println(model)
	// Create a context for the API call.
	ctx := context.Background()

	// Create a new VertexAI client.
	client, err := genai.NewClient(ctx, project, region)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close() // Close the client when the function exits.

	gemini := client.GenerativeModel(model)
	prompt := genai.Text(chatPrompt)

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}

	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	// // fmt.Println(string(rb))

	// Unmarshal the JSON response into a struct.
	var response struct {
		Candidates []struct {
			Content struct {
				Parts []json.RawMessage `json:"Parts"`
			} `json:"Content"`
		} `json:"Candidates"`
	}

	if err := json.Unmarshal(rb, &response); err != nil {
		log.Printf("Error unmarshaling JSON response: %v", err)
		return "", fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	// Check if the response contains any candidates.
	if len(response.Candidates) == 0 {
		log.Println("No candidates found in the response")
		return "", fmt.Errorf("no candidates found in the response")
	}

	// Decode the response parts into a string slice.
	var parts []string
	for _, rawPart := range response.Candidates[0].Content.Parts {
		var textPart string
		if err := json.Unmarshal(rawPart, &textPart); err != nil {
			return "", fmt.Errorf("error unmarshaling part: %v", err)
		}
		parts = append(parts, textPart)
	}

	// fmt.Println(parts)

	// Join the decoded string parts into a single string
	// Validate the JSON before returning
	// validatedJSON, err := utils.ValidateJSON(strings.Join(parts, ""))
	// if err != nil {
	// 	return "", fmt.Errorf("error validating JSON: %w", err)
	// }
	// return validatedJSON, nil

	return strings.Join(parts, ""), nil

}
