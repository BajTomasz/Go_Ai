package Tasks

import (
	"Go_Ai/APIs"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func Liar() {
	taskToken, _, secrets := APIs.DownloadTask("liar")

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
	APIs.CheckResponse(respForm, err)

	var task Task
	json.NewDecoder(respForm.Body).Decode(&task)

	var messages []APIs.Message
	messages = append(messages, APIs.Message{
		Role:    "system",
		Content: "You are a lie detector, you will get a question and then an answer. Return YES if the answer is true. If the answer is a lie, return NO.",
	})
	messages = append(messages, APIs.Message{
		Role:    "user",
		Content: question + "/\n###/\n" + task.Answer,
	})

	response := APIs.Completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, 0, nil)
	result := response.Choices[0].Message.Content
	fmt.Println(messages)
	fmt.Println(response)

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
