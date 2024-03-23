package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("helloapi")

	//____Solve_Task____
	type Task struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Cookie string `json:"cookie"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println("Code:", task.Code)
	fmt.Println("Msg:", task.Msg)
	fmt.Println("Cookie:", task.Cookie)

	postBody, _ := json.Marshal(map[string]string{
		"answer": task.Cookie,
	})
	sendAnswer(taskToken, postBody, secrets)
}
