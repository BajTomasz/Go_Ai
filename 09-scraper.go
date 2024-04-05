package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func scraper() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("scraper")

	//____Solve_Task____
	type Task struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Input    string `json:"input"`
		Question string `json:"question"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)

	client := &http.Client{}
	req, err := http.NewRequest("GET", task.Input, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	var textResp *http.Response
	var bytes []byte
	for {
		fmt.Println("Send GET to " + task.Input)
		textResp, err = client.Do(req)
		if textResp.StatusCode == 200 {
			defer textResp.Body.Close()
			bytes, err = io.ReadAll(textResp.Body)
			checkError(err)
			break
		} else {
			fmt.Printf("ERROR: %v\n", err)
			fmt.Println("Status HTTP:", textResp.Status)
			b, _ := io.ReadAll(textResp.Body)
			fmt.Printf("%s", b)
			time.Sleep(3 * time.Second)
		}

	}
	checkResponse(textResp, err)

	fmt.Println(string(bytes))
	checkError(err)

	fmt.Println(task)

	var messages []Message
	messages = append(messages, Message{
		Role:    "system",
		Content: string(bytes) + "/\n" + "###" + "/\n" + task.Msg,
	})
	messages = append(messages, Message{
		Role:    "user",
		Content: task.Question,
	})

	var max_tokens int64
	max_tokens = 100
	response := completions(secrets.OpenaiAPIKey, "gpt-3.5-turbo-0125", messages, max_tokens, nil)
	result := response.Choices[0].Message.Content
	fmt.Println(response)

	postBody, _ := json.Marshal(map[string]string{
		"answer": result,
	})
	sendAnswer(taskToken, postBody, secrets)
}
