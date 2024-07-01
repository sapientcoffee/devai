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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kyokomi/emoji/v2"
)

type ReviewItem struct {
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Description string `json:"description"`
}

func ReadFileContents(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func BuildChatPrompt(basePrompt, promptFlag, code string) string {
	if promptFlag != "" {
		fmt.Println("Additial prompt flag set: " + promptFlag)
		basePrompt += " and " + promptFlag
	}
	return basePrompt + " for the following code: " + code
}

// ValidateJSON takes a JSON string and attempts to fix common issues.
// It returns the corrected JSON string and any errors encountered.
func ValidateJSON(jsonString string) (string, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		// If the JSON is invalid, try to fix it using FixJSON
		fixedJSON, err := FixJSON(jsonString)
		if err != nil {
			return "", fmt.Errorf("invalid JSON: %w", err)
		}
		return fixedJSON, nil
	}
	// Marshal the corrected data back to JSON
	correctedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling corrected JSON: %w", err)
	}
	return string(correctedJSON), nil
}

// FixJSON attempts to fix common JSON issues and return valid JSON.
// It returns the corrected JSON string and any errors encountered.
func FixJSON(output string) (string, error) {
	fmt.Println("Attempting to fix JSON syntax issues...")

	// 1. Preliminary Cleanup:
	// Remove leading/trailing whitespace and backticks (common in shell output):
	output = strings.TrimSpace(output)
	output = strings.Trim(output, "`")

	// Remove the word "json" if it's at the beginning (common in some APIs):
	output = strings.TrimPrefix(output, "json")

	// 2. Robust JSON Parsing with Error Handling:
	var parsedData interface{}
	decoder := json.NewDecoder(strings.NewReader(output))

	// Configure the decoder to be more lenient about certain deviations:
	decoder.UseNumber()             // Keep numbers as JSON numbers, not convert to float64
	decoder.DisallowUnknownFields() // Comment out if you expect extra fields in JSON

	err := decoder.Decode(&parsedData)
	if err == nil {
		// If parsing succeeds, we have valid JSON; just reformat it:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ") // Optional: pretty-print with indentation
		if err := encoder.Encode(parsedData); err != nil {
			return "", fmt.Errorf("error re-encoding JSON: %v", err)
		}
		return buf.String(), nil
	}

	// 3. Handle Common JSON Syntax Errors:
	// If initial parsing failed, try common fixes:
	if strings.HasPrefix(err.Error(), "invalid character") && strings.Contains(err.Error(), "after top-level value") {
		// If the error indicates multiple JSON objects, try to split and parse them individually:
		var fixedOutput strings.Builder
		for _, part := range strings.Split(output, "\n") { // Split on newlines
			part = strings.TrimSpace(part)
			if part != "" {
				_, err := FixJSON(part) // Recursively fix each part
				if err != nil {
					return "", err // Return the original error if fixing fails
				}
				fixedOutput.WriteString(part)
				fixedOutput.WriteString("\n") // Add newline between objects
			}
		}
		return fixedOutput.String(), nil

	} else if err == io.EOF {
		// Handle unexpected end of JSON input
		output = output + "}"
		return FixJSON(output)
	} else {
		// If it's another error type, return for further investigation
		return "", fmt.Errorf("JSON parsing error: %v", err)
	}
}

func CreateTable(reviewItems []ReviewItem, format string) string {
	var tableBuilder strings.Builder

	switch format {
	case "cli":
		// CLI Output using Emoji
		t := table.NewWriter()
        t.SetOutputMirror(os.Stdout)
        t.AppendHeader(table.Row{"Question", "Answer", "Description"})


		for _, item := range reviewItems {
			var answerDisplay string
			if item.Answer == "yes" {
				answerDisplay = emoji.Sprint(":white_check_mark:")
			} else if item.Answer == "no" {
				answerDisplay = emoji.Sprint(":x:")
			} else {
				answerDisplay = item.Answer
			}

			t.AppendRow(table.Row{item.Question, answerDisplay, item.Description})
		}

		t.Render()
		return "" // No string output needed for CLI

	case "markdown":
		// Markdown Output
		tableBuilder.WriteString("| Question    | Answer | Description |\n")
		tableBuilder.WriteString("|------------|--------|-------------|\n")

		for _, item := range reviewItems {
			tableBuilder.WriteString(fmt.Sprintf("| %s | %s | %s |\n", item.Question, item.Answer, item.Description))
		}

		return tableBuilder.String()

	default:
		// Handle invalid format gracefully
		return "Invalid format specified. Please choose 'cli' or 'markdown'."
	}
}
