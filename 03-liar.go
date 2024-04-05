package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func liar() {
	taskToken, _, secrets := downloadTask("liar")

	//____Solve_Task____
	type Task struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Answer string `json:"answer"`
	}

	question := "What are traditional Polish foods?"
	formData := url.Values{
		"question": {question},
	}

	urlQuestion := secrets.Url + "task/" + taskToken
	fmt.Println("Send post to " + urlQuestion)
	respForm, err := http.PostForm(urlQuestion, formData)
	checkResponse(respForm, err)

	var task Task
	json.NewDecoder(respForm.Body).Decode(&task)

	var messages []Message
	messages = append(messages, Message{
		Role:    "system",
		Content: "You are a lie detector, you will get a question and then an answer. Return YES if the answer is true. If the answer is a lie, return NO.",
	})
	messages = append(messages, Message{
		Role:    "user",
		Content: question + "/\n###/\n" + task.Answer,
	})

	response := completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, 0, nil)
	result := response.Choices[0].Message.Content
	fmt.Println(messages)
	fmt.Println(response)

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	sendAnswer(taskToken, postBody, secrets)
}
