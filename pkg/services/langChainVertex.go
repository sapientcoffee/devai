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

// This file contains the code for the Chat function, which uses the
// Google Cloud VertexAI Codey API to generate text responses to user prompts.

package genaiService

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
)

func LangChainVertexChat(prompt string) string {
	ctx := context.Background()
	// project := os.Getenv("VERTEX_PROJECT")
	project := "coffee-and-codey"
	// location := os.Getenv("VERTEX_LOCATION")
	region := "us-central1"
	llm, err := vertex.New(
		ctx,
		googleai.WithCloudProject(project),
		googleai.WithCloudLocation(region),
		// googleai.WithCredentialsFile(credentialsJSONFile),
	)
	if err != nil {
		log.Fatal(err)
	}

	// prompt := "What is the L2 Lagrange point?"
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(answer)
	return answer
}
