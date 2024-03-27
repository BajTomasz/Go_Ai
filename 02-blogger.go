package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func blogger() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("blogger")

	//____Solve_Task____
	type Task struct {
		Code int      `json:"code"`
		Msg  string   `json:"msg"`
		Task []string `json:"blog"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println(task)

	var messages []Message
	var results []string

	for _, chapter := range task.Task {
		messages = []Message{}
		messages = append(messages, Message{
			Role:    "system",
			Content: task.Msg,
		})
		messages = append(messages, Message{
			Role:    "user",
			Content: chapter,
		})
		response := completion(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages)
		results = append(results, response.Choices[0].Message.Content)
		fmt.Println(messages)
		fmt.Println(response)

	}

	postBody, _ := json.Marshal(map[string][]string{
		"answer": results,
	})
	sendAnswer(taskToken, postBody, secrets)
}
