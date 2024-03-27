package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func moderation() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("moderation")

	//____Solve_Task____
	type Task struct {
		Code  int      `json:"code"`
		Msg   string   `json:"msg"`
		Input []string `json:"input"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println(task)

	var results []int
	for i := range task.Input[:] {
		if moderations(secrets.OpenaiAPIKey, task.Input[i]).Results[0].Flagged {
			results = append(results, 1)
		} else {
			results = append(results, 0)
		}
	}

	postBody, _ := json.Marshal(map[string][]int{
		"answer": results,
	})
	sendAnswer(taskToken, postBody, secrets)
}
