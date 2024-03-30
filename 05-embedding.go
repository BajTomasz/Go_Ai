package main

import (
	"bytes"
	"encoding/json"
)

func embedding() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("embedding")

	//____Solve_Task____
	type Task struct {
		Code     int      `json:"code"`
		Msg      string   `json:"msg"`
		Input    []string `json:"input"`
		Question string   `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)

	textToEmbed := []string{"Hawaiian pizza"}
	response := embeddings(secrets.OpenaiAPIKey, "text-embedding-ada-002", textToEmbed)
	result := response.Data[0].Embedding

	postBody, _ := json.Marshal(map[string][]float32{
		"answer": result,
	})
	sendAnswer(taskToken, postBody, secrets)
}
