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

// This file contains the code for the Chat function, which uses the
// Google Cloud VertexAI Codey API to generate text responses to user prompts.
package genaiService

import "fmt"

// ChatEngine represents the different AI chat engines available.
// This enum allows users to select the desired engine for text generation.
type ChatEngine string

const (
	// LangChainEngine represents the LangChain chat engine, which is known for its flexibility and integration with various LLM models.
	LangChainEngine ChatEngine = "LangChain"
	// CodeyEngine represents the Codey chat engine, a Google Cloud service specifically designed for code generation and related tasks.
	CodeyEngine ChatEngine = "codey"
	// VertexEngine represents the Vertex AI (via its SDK) chat engine, a powerful and scalable platform for deploying and managing machine learning models.
	VertexEngine ChatEngine = "vertex"
)

// DefaultChatEngine is the chat engine used if none is specified.
// VertexEngine is the default as it offers a robust and well-supported platform for AI tasks.
var DefaultChatEngine ChatEngine = VertexEngine

// GetAIResponse interacts with the specified AI service to generate a text response based on the provided prompt.
//
// Args:
//
//	chatPrompt: The user's input prompt for the AI to respond to.
//	chatEngine: The desired AI chat engine to use. If not specified, the default engine (VertexEngine) will be used.
//
// Returns:
//
//	response: The generated text response from the AI service.
//	error: Any error encountered during the interaction with the AI service.
func GetAIResponse(chatPrompt string, chatEngine ChatEngine) (string, error) {
	var response string
	var err error

	// If no chat engine is specified, use the default.
	if chatEngine == "" {
		chatEngine = DefaultChatEngine
	}

	switch chatEngine {
	case LangChainEngine:
		fmt.Printf(".. Using Langchain with the model: ")
		response = LangChainVertexChat(chatPrompt)
		if response == "" {
			return "", fmt.Errorf("error calling LangChainVertexChat: %w", err)
		}
	case CodeyEngine:
		fmt.Printf(".. Using Codey LLM")
		response = CodeyChat(chatPrompt)
		if response == "" {
			return "", fmt.Errorf("error calling CodeyChat: %w", err)
		}
	case VertexEngine:
		fmt.Printf(".. Using VertexAI SDK with the model: ")
		response, err = VertexChatText(chatPrompt)
		if err != nil {
			return "", fmt.Errorf("error calling VertexChatText: %w", err)
		}
	default:
		return "", fmt.Errorf("invalid chat engine: %s", chatEngine)

	}

	return response, nil
}
