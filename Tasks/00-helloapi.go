package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
)

func Helloapi() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("helloapi")
	//fmt.Println(resp.String())

	//____Solve_Task____
	type Task struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Cookie string `json:"cookie"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)
	fmt.Println(task)

	postBody, _ := json.Marshal(map[string]string{
		"answer": task.Cookie,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
