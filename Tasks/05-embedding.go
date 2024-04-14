package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
)

func Embedding() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("embedding")

	//____Solve_Task____
	type Task struct {
		Code     int      `json:"code"`
		Msg      string   `json:"msg"`
		Input    []string `json:"input"`
		Question string   `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	textToEmbed := []string{"Hawaiian pizza"}
	response := APIs.Embeddings(secrets.OpenaiAPIKey, "text-embedding-3-small", textToEmbed)
	result := response.Data[0].Embedding

	postBody, _ := json.Marshal(map[string][]float32{
		"answer": result,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
