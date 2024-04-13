package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
)

func Blogger() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("blogger")

	//____Solve_Task____
	type Task struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Task []string `json:"blog"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)
	fmt.Println(task)

	var messages []APIs.Message
	var results []string

	for _, chapter := range task.Task {
		messages = []APIs.Message{}
		messages = append(messages, APIs.Message{
			Role:    "system",
			Content: task.Msg,
		})
		messages = append(messages, APIs.Message{
			Role:    "user",
			Content: chapter,
		})
		response := APIs.Completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, 0, nil)
		results = append(results, response.Choices[0].Message.Content)
		fmt.Println(messages)
		fmt.Println(response)

	}

	postBody, _ := json.Marshal(map[string][]string{
		"answer": results,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
