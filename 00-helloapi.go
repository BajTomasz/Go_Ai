package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func helloapi() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("helloapi")
	//fmt.Println(resp.String())

	//____Solve_Task____
	type Task struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Cookie string `json:"cookie"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println(task)

	postBody, _ := json.Marshal(map[string]string{
		"answer": task.Cookie,
	})
	sendAnswer(taskToken, postBody, secrets)
}
