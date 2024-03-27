package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ModerationRequest struct {
	Input string `json:"input"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Messages   []Message `json:"messages"`
	Model      string    `json:"model"`
	Max_tokens int64     `json:"max_tokens,omitempty"`
	N          int64     `json:"n,omitempty"`
}

type Choice struct {
	Index         int     `json:"index"`
	Message       Message `json:"message"`
	Logprobs      int     `json:"logprobs"`
	Finish_reason string  `json:"finish_reason"`
}

type Usage struct {
	Prompt_tokens     int `json:"prompt_tokens"`
	Completion_tokens int `json:"completion_tokens"`
	Total_tokens      int `json:"total_tokens"`
}
type CompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

func moderations(openaiAPIKey string, input string) Moderation {
	//start := time.Now()

	moderationRequest := ModerationRequest{
		Input: input,
	}
	jsonInput, _ := json.Marshal(moderationRequest)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/moderations", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	respModer, err := http.DefaultClient.Do(req)
	checkResponse(respModer, err)
	defer respModer.Body.Close()

	var moderation Moderation
	json.NewDecoder(respModer.Body).Decode(&moderation)
	//fmt.Println(moderation)
	//fmt.Println("%.2fs", time.Since(start).Seconds())
	return moderation
}

func completion(openaiAPIKey string, model string, messages []Message) CompletionResponse {
	url := "https://api.openai.com/v1/chat/completions"

	request := CompletionRequest{
		Model:    model,
		Messages: messages,
	}

	jsonInput, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	response, err := http.DefaultClient.Do(req)
	checkResponse(response, err)
	defer response.Body.Close()

	var result CompletionResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	checkError(err)
	return result
}
