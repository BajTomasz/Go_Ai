package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
	"fmt"
)

func Moderation() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("moderation")

	//____Solve_Task____
	type Task struct {
		Code  int      `json:"code"`
		Msg   string   `json:"msg"`
		Input []string `json:"input"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)
	fmt.Println(task)

	var results []int
	for i := range task.Input[:] {
		if APIs.Moderations(secrets.OpenaiAPIKey, task.Input[i]).Results[0].Flagged {
			results = append(results, 1)
		} else {
			results = append(results, 0)
		}
	}

	postBody, _ := json.Marshal(map[string][]int{
		"answer": results,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
