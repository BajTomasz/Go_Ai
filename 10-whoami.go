package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Task struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Hint string `json:"hint"`
}

type AkinatorResponse struct {
	Sure string `json:"sure"`
	Name string `json:"name"`
}

func downloadWhoami() (string, Task, Secrets) {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("whoami")

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	return taskToken, task, secrets
}

func whoami() {
	var informations, taskToken string
	var task Task
	var secrets Secrets
	var messages []Message
	var response CompletionResponse
	var akinatorResponse AkinatorResponse

	// ____Solve_Task____
	messages = append(messages, Message{
		Role:    "system",
		Content: "Based on the clues, guess who I'm talking about. Answer \"YES\" if you are really sure who I am. If you need more instructions, write \"NO\". Reply in JSON format: {\"sure\": \"YES\", \"name\": \"Ben Smith\"}",
	})
	messages = append(messages, Message{
		Role:    "user",
		Content: informations,
	})

	count := 3
	for i := 0; i <= count; i++ {
		taskToken, task, secrets = downloadWhoami()
		informations = informations + "\n" + "###" + "\n" + task.Hint
		messages[1].Content = informations
		fmt.Println(i, messages[1].Content)
		if i == count {
			response = completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, 0, nil)
			fmt.Println(response)
			err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &akinatorResponse)
			checkError(err)
			if akinatorResponse.Sure == "YES" {
				break
			} else {
				count = count + 1
			}
		}
		if i > 6 {
			break
		}
	}

	postBody, _ := json.Marshal(map[string]string{
		"answer": akinatorResponse.Name,
	})
	sendAnswer(taskToken, postBody, secrets)
}
