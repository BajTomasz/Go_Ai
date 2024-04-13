package Tasks

import (
	"Go_Ai/APIs"
	"bytes"
	"encoding/json"
)

func Functions() {
	var resp bytes.Buffer
	taskToken, resp, secrets := APIs.DownloadTask("functions")

	//____Solve_Task____
	type Task struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	APIs.CheckError(err)

	parameters := APIs.Parameters{
		Type: "object",
		Properties: map[string]APIs.Property{
			"name": {
				Type:        "string",
				Description: "User name",
			},
			"surname": {
				Type:        "string",
				Description: "User Surname",
			},
			"year": {
				Type:        "integer",
				Description: "User's age",
			},
		},
	}

	funcObj := APIs.Function{
		Name:        "addUser",
		Description: "Send me definition of function named addUser that require 3 params",
		Parameters:  parameters,
	}

	postBody, _ := json.Marshal(map[string]APIs.Function{
		"answer": funcObj,
	})
	APIs.SendAnswer(taskToken, postBody, secrets)
}
