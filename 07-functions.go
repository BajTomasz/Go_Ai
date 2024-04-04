package main

import (
	"bytes"
	"encoding/json"
)

func functions() {
	var resp bytes.Buffer
	taskToken, resp, secrets := downloadTask("functions")

	//____Solve_Task____
	type Task struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Hint string `json:"hint"`
	}

	var task Task
	err := json.NewDecoder(&resp).Decode(&task)
	checkError(err)

	parameters := Parameters{
		Type: "object",
		Properties: map[string]Property{
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

	funcObj := Function{
		Name:        "addUser",
		Description: "Send me definition of function named addUser that require 3 params",
		Parameters:  parameters,
	}

	postBody, _ := json.Marshal(map[string]Function{
		"answer": funcObj,
	})
	sendAnswer(taskToken, postBody, secrets)
}
