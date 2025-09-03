package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type LLM_RESPONSE struct {
	Urls      []string `json:"urls"`
	FoundData []string `json:"foundData"` //interface{}
}

func callLLM(responseBody []byte, prompt string) (LLM_RESPONSE, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	var llmResponse LLM_RESPONSE

	api_key := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(option.WithAPIKey(api_key))
	context := context.Background()
	chatCompletion, err := client.Chat.Completions.New(context, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.AssistantMessage(getWebCrawlerSystemPrompt(prompt)),
			openai.UserMessage(string(responseBody))},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		fmt.Println("Error creating chat completion:", err)
		return llmResponse, err
	}
	content := strip(chatCompletion.Choices[0].Message.Content)
	// fmt.Println("Content Response:", content)
	var temp struct {
		Urls      []string    `json:"urls"`
		FoundData interface{} `json:"foundData"` //interface{}
	}
	err = json.Unmarshal([]byte(content), &temp)

	llmResponse.Urls = temp.Urls
	llmResponse.FoundData = interfaceArrayToStringSlice(temp.FoundData)

	if err != nil {
		fmt.Println("Error parsing LLM response:", err)
		return llmResponse, err
	}
	// fmt.Println("Chat Completion:", chatCompletion.Choices[0].Message.Content)
	return llmResponse, nil
}

func getUserAnswer(foundData []string, prompt string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	api_key := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(option.WithAPIKey(api_key))
	context := context.Background()
	chatCompletion, err := client.Chat.Completions.New(context, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{openai.AssistantMessage(getUserAnswerSystemPrompt(prompt)),
			openai.UserMessage(strings.Join(foundData, "\n"))},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		fmt.Println("Error creating chat completion:", err)
		return ""
	}
	content := chatCompletion.Choices[0].Message.Content
	return content
}
