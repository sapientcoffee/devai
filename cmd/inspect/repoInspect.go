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
package inspect

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
	// dir stores the directory to be inspected.
	dir string
	// promptFlag stores the prompt to be used for the AI.
	promptFlag string
)


// repoInspect inspects the given directory and generates reports for each file.
func repoInspect(dir string) {
	fmt.Println("Generating reports...")

	// Create the output directory if it doesn't exist.
	outputDir := "inspect-output"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Create a file for the directory tree.
	archFilePath := filepath.Join(outputDir, "repo-arch.txt")
	archFile, err := os.Create(archFilePath)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", archFilePath, err)
		return
	}
	defer archFile.Close()

	// Walk the directory and generate reports.
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude hidden directories, except the root directory.
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != dir {
			return filepath.SkipDir
		}

		// Generate the directory tree structure.
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		treePrefix := generateTreePrefix(relPath, info.IsDir())
		if !info.IsDir() || !strings.HasPrefix(info.Name(), ".") {
			_, err = archFile.WriteString(treePrefix + info.Name() + "\n")
			if err != nil {
				return err
			}
		}

		// Process files with supported extensions (excluding the directory tree file).
        if !info.IsDir() && isSupportedFile(info.Name()) && info.Name() != "repo-arch.txt" {
           fmt.Printf("Processing file: %s\n", path)

			// Read the file content.
			fileContent, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			// Build the prompt for the AI.
			chatPrompt := utils.BuildChatPrompt(
				`### instruction ###
				You are an experinced software architectect and developer, please explain the code to another experinced developer. The purpose is to help undertsnand what the code is doing and any business logic to help increase understanding of the code base.


				`,
				promptFlag, string(fileContent))

			// Get AI response.
			response, err := genaiService.GetAIResponse(chatPrompt, "")
			if err != nil {
				fmt.Printf("Error getting AI response: %v\n", err)
				return err
			}

			// Add file name and full path to the response.
			responseWithMetadata := fmt.Sprintf("# %s\n_path: %s_\n\n%s", info.Name(), path, response)

			// Write the response to a file.
			outputFileName := info.Name() + ".md"
			outputFilePath := filepath.Join(outputDir, outputFileName)
			if err := os.WriteFile(outputFilePath, []byte(responseWithMetadata), 0644); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				return err
			}

			fmt.Printf("Output written to: %s\n", outputFilePath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	fmt.Println("Reports generated successfully!")
}

// generateTreePrefix generates the prefix for the directory tree structure.
func generateTreePrefix(relPath string, isDir bool) string {
	prefix := ""
	depth := len(strings.Split(relPath, string(os.PathSeparator))) - 1
	if depth > 0 { // Add indentation for subdirectories
		for i := 0; i < depth-1; i++ {
			prefix += "│   "
		}
		if isDir {
			prefix += "├── "
		} else {
			prefix += "│   ├── " // Indentation for files
		}
	} else if depth == 0 && !isDir { // Handle files in the root directory
		prefix += "├── "
	}
	return prefix
}

// isSupportedFile checks if the file has a supported extension.
func isSupportedFile(fileName string) bool {
    supportedExtensions := []string{".go", ".py", ".js", ".yaml", ".yml", ".tf", ".cpp", ".cxx", ".cc", ".h", ".hpp"} // Add more extensions as needed
    for _, ext := range supportedExtensions {
        if strings.HasSuffix(fileName, ext) {
            return true
		}
		}
	return false
}

// repoInspectCmd represents the repoInspect command.
var repoInspectCmd = &cobra.Command{
	Use:   "repo",
	Short: "Inspect a repository and generate reports for each file.",
	Long:  `This command inspects a repository and generates reports for each file using an AI.`,
	Run: func(cmd *cobra.Command, args []string) {
		dirFlag, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file flag:", err)
			return
		}

		if dirFlag != "" {
			dir = dirFlag
		} else {
			fmt.Println("No directory specified, using current directory.")
			dir = "." // Default to the current directory
		}
		fmt.Println("Inspecting repository:", dirFlag)
		repoInspect(dirFlag)
	},
}

func init() {
	repoInspectCmd.Flags().StringP("file", "f", "", "The directory containing the code to be reviewed")
}