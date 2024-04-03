package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type Usage struct {
	Prompt_tokens     int `json:"prompt_tokens"`
	Completion_tokens int `json:"completion_tokens"`
	Total_tokens      int `json:"total_tokens"`
}

// Moderations
type ModerationRequest struct {
	Input string `json:"input"`
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

// Completions
type CompletionRequest struct {
	Messages   []Message `json:"messages"`
	Model      string    `json:"model"`
	Max_tokens int64     `json:"max_tokens,omitempty"`
	N          int64     `json:"n,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index         int     `json:"index"`
	Message       Message `json:"message"`
	Logprobs      int     `json:"logprobs"`
	Finish_reason string  `json:"finish_reason"`
}

type CompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

func completions(openaiAPIKey string, model string, messages []Message) CompletionResponse {
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

// Embeddings
type EmbeddingRequest struct {
	Input           []string `json:"input"`
	Model           string   `json:"model"`
	Encoding_format string   `json:"encoding_format,omitempty"`
}

type EmbeddingObject struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
}

type EmbeddingResponse struct {
	Object string            `json:"object"`
	Data   []EmbeddingObject `json:"data"`
	Model  string            `json:"model"`
	Usage  Usage             `json:"usage"`
}

func embeddings(openaiAPIKey string, model string, input []string) EmbeddingResponse {
	url := "https://api.openai.com/v1/embeddings"

	embeddingRequest := EmbeddingRequest{
		Input: input,
		Model: model,
	}

	jsonInput, _ := json.Marshal(embeddingRequest)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	response, err := http.DefaultClient.Do(req)
	checkResponse(response, err)
	defer response.Body.Close()

	var embeddingResponse EmbeddingResponse
	err = json.NewDecoder(response.Body).Decode(&embeddingResponse)
	checkError(err)
	return embeddingResponse
}

// Transcriptions
type TranscriptionResponse struct {
	Teskt string `json:"text"`
}

func transcriptions(openaiAPIKey string, model string, audioReader []byte) TranscriptionResponse {
	transcriptionsUrl := "https://api.openai.com/v1/audio/transcriptions"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("model", model)
	writer.WriteField("language", "pl")

	part, err := writer.CreateFormFile("file", "audio.mp3")
	part.Write(audioReader)
	writer.Close()

	req, err := http.NewRequest("POST", transcriptionsUrl, &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	response, err := http.DefaultClient.Do(req)
	checkResponse(response, err)
	defer response.Body.Close()
	var buf bytes.Buffer
	io.Copy(&buf, response.Body)

	var transcriptionResponse TranscriptionResponse
	err = json.NewDecoder(&buf).Decode(&transcriptionResponse)

	checkError(err)
	return transcriptionResponse
}
