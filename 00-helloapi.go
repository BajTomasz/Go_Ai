package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	var resp bytes.Buffer
	resp = downloadTask()

	//____Solve_Task____
	type Task struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Cookie string `json:"cookie"`
	}

	var task Task
	//err := json.NewDecoder(resp.Body).Decode(&task)
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)
	fmt.Println("Code:", task.Code)
	fmt.Println("Msg:", task.Msg)
	fmt.Println("Token:", task.Cookie)

	answer := strings.NewReader("{\"answer\":\"" + task.Cookie + "\"}")
	sendAnswer(answer)
}
